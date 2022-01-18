package log

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	rootLogger atomic.Value
)

const (
	LabelErr = "error"
)

//RootLogger is the structure to use to produce logs.
type RootLogger struct {
	// Used lowLogger
	lowLogger *zap.Logger
	//log config
	ShowDebug bool
	//log context part
	component  string
	appVersion string
}

//InitLogger create and initialize an atomic value (used as singleton) of the logs layer
func InitLogger(component string, debug bool, version string) RootLogger {
	logger := RootLogger{
		ShowDebug:  debug,
		component:  component,
		appVersion: version,
	}

	// Min log level
	dyn := zap.NewAtomicLevel()
	dyn.SetLevel(zap.InfoLevel)
	if logger.ShowDebug {
		dyn.SetLevel(zap.DebugLevel)
	}

	logger.lowLogger, _ = zap.Config{
		Level:             dyn,
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
	return GetLogger()
}

//Allow to encode duration in ms
func MilliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(int64(d) / int64(time.Millisecond))
}

func GetLogger() RootLogger {
	return rootLogger.Load().(RootLogger)
}

/*
Debug is for all logs that are for developers.
*/
func Debug(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = addContextInfo(ctx, fields...)
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
	return allFields
}

// Log a panic content
func LogPanic() {
	if r := recover(); r != nil {
		err := fmt.Errorf("%v", r)
		Error(context.Background(), "Panic", err, zap.Stack("stack"))
	}
}
