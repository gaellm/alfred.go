/*
 * Copyright The Alfred.go Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package log

import (
	"alfred/internal/tracing"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	rootLogger atomic.Value
)

const (
	LabelErr = "error"
)

const LOG_LEVEL_DEBUG = "debug"
const LOG_LEVEL_INFO = "info"
const LOG_LEVEL_ERROR = "error"

var LOG_LEVEL_LIST = []string{LOG_LEVEL_DEBUG, LOG_LEVEL_ERROR, LOG_LEVEL_INFO}
var ZAP_LEVEL_MAP = map[string]zapcore.Level{
	LOG_LEVEL_INFO:  zap.InfoLevel,
	LOG_LEVEL_DEBUG: zap.DebugLevel,
	LOG_LEVEL_ERROR: zap.ErrorLevel,
}

//RootLogger is the structure to use to produce logs.
type RootLogger struct {
	// Used lowLogger
	lowLogger *zap.Logger
	//log config
	ShowDebug bool
	//log context part
	component  string
	appVersion string
	dyn        zap.AtomicLevel
}

//InitLogger create and initialize an atomic value (used as singleton) of the logs layer
func InitLogger(component string, debug bool, version string) RootLogger {
	logger := RootLogger{
		ShowDebug:  debug,
		component:  component,
		appVersion: version,
	}

	// Min log level
	logger.dyn = zap.NewAtomicLevel()
	logger.dyn.SetLevel(zap.InfoLevel)
	if logger.ShowDebug {
		logger.dyn.SetLevel(zap.DebugLevel)
	}

	logger.lowLogger, _ = zap.Config{
		Level:             logger.dyn,
		Development:       false,
		Encoding:          "json",
		DisableCaller:     true,
		DisableStacktrace: true,
		EncoderConfig: zapcore.EncoderConfig{ //used names and encoder for common fields
			TimeKey:       "dateTime",
			LevelKey:      "level",
			NameKey:       "lowLogger",
			MessageKey:    "msg",
			StacktraceKey: "stacktrace",
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			/* si on préfère le format normalisé RFC3339 :
			EncodeTime:     func (t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(time.RFC3339))
			},*/
			EncodeDuration: zapcore.StringDurationEncoder, //String() function will be used to log a duration
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		InitialFields: map[string]interface{}{ //Additional mandatory fields
			"component":  logger.component,
			"appVersion": logger.appVersion,
		},
	}.Build()

	rootLogger.Store(logger)

	tracingErrorHandler()

	return GetLogger()
}

//Allow to encode duration in ms
func MilliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(int64(d) / int64(time.Millisecond))
}

func GetLogger() RootLogger {
	return rootLogger.Load().(RootLogger)
}

func addTracingContextFields(ctx context.Context, fields ...zapcore.Field) []zapcore.Field {

	traceId, spanId, flags := tracing.GetSpanContext(tracing.GetSpanFromContext(ctx))
	if traceId == "00000000000000000000000000000000" {
		return fields
	}

	tracingContextFields := []zapcore.Field{
		zap.String("trace-id", traceId),
		zap.String("span-id", spanId),
		zap.String("trace-flags", flags),
	}

	return append(fields, tracingContextFields...)

}

/*
Debug is for all logs that are for developers.
*/
func Debug(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = addContextInfo(ctx, fields...)
	fields = appendLogFields(ctx, fields...)
	GetLogger().lowLogger.Debug(msg, fields...)
}

/*
Info is for all logs that are for informative purpose.
Incoming requests for example, or outgoing results.
*/
func Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	allFields := appendLogFields(ctx, fields...)

	GetLogger().lowLogger.Info(msg, allFields...)
}

/*
Change level at runtime
*/
func SetLevel(levelStr string) error {

	for _, level := range LOG_LEVEL_LIST {
		if strings.EqualFold(levelStr, level) {

			GetLogger().dyn.SetLevel(ZAP_LEVEL_MAP[level])
			return nil
		}
	}

	return errors.New("Logger level " + levelStr + " does not exist. It should be one of " + fmt.Sprint(LOG_LEVEL_LIST))

}

/*
Get current level at runtime
*/
func GetLevel() string {

	return GetLogger().dyn.Level().String()
}

/*
Info is for all logs that can help to investigate a specific behavior but doesn't indicate an issue that needs to be investigated
Incoming requests for example, or outgoing results.
*/
func Warn(ctx context.Context, msg string, err error, fields ...zapcore.Field) {
	allFields := appendLogFields(ctx, fields...)

	errMsg := "nil"
	if err != nil {
		errMsg = err.Error()
	}
	allFields = append(allFields, zap.String(LabelErr, errMsg))

	GetLogger().lowLogger.Warn(msg, allFields...)
}

/*
Error level is for any error that is not recoverable: a retry will not solve the problem.
The administrator should have action to fix the problem.
It is for inconsistency in data for example.
*/
func Error(ctx context.Context, msg string, err error, fields ...zapcore.Field) {
	allFields := appendLogFields(ctx, fields...)

	errMsg := "nil"
	if err != nil {
		errMsg = err.Error()
	}
	allFields = append(allFields, zap.String(LabelErr, errMsg))

	span := tracing.GetSpanFromContext(ctx)
	tracing.SetSpanStatusError(&span, err)

	GetLogger().lowLogger.Error(msg, allFields...)
}

/*
Fatal is for recoverable errors. Executable should be stopped after this error.
It is intended to cover cases like network error, database overload, lack of resources, etc.
*/
func Fatal(ctx context.Context, msg string, err error, fields ...zapcore.Field) {
	allFields := appendLogFields(ctx, fields...)

	errMsg := "nil"
	if err != nil {
		errMsg = err.Error()
	}
	allFields = append(allFields, zap.String(LabelErr, errMsg))

	GetLogger().lowLogger.Fatal(msg, allFields...)
}

//Add fields from ctx
func addContextInfo(ctx context.Context, fields ...zapcore.Field) []zapcore.Field {
	if ctx != nil {

		//---- RequestId : identify an input request between different services
		requestId := ctx.Value("X-RequestId")
		if requestId != nil {
			fields = append(fields, zap.String("RequestId", requestId.(string)))
		}

	}
	return fields
}

func appendLogFields(ctx context.Context, fields ...zapcore.Field) []zapcore.Field {
	allFields := []zapcore.Field{}
	allFields = addContextInfo(ctx, allFields...)
	allFields = append(allFields, fields...)
	allFields = addTracingContextFields(ctx, allFields...)
	return allFields
}

// Log a panic content
func LogPanic() {
	if r := recover(); r != nil {
		err := fmt.Errorf("%v", r)
		Error(context.Background(), "Panic", err, zap.Stack("stack"))
	}
}

func tracingErrorHandler() {

	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		Error(context.Background(), "otel error", err)
	}))
}
