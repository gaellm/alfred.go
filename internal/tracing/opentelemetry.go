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

package tracing

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/imdario/mergo"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	TracerName    = "instrumentation/alfred.go"
	TracerVersion = "1.0"
)

type OtelConfig struct {
	ServiceName           string
	ServiceVersion        string
	ServiceNamespace      string
	DeploymentEnvironment string
	ExporterOtlpEndpoint  string
	ExporterInsecure      bool
	TracesSampler         string
	TracesSamplerArg      string
}

var DEFAULT_OTEL_CONFIG = OtelConfig{
	ServiceName:           "unknown_service",
	ServiceNamespace:      "default",
	DeploymentEnvironment: "all",
	ExporterInsecure:      true,
	TracesSampler:         "parentbased_traceidratio",
	TracesSamplerArg:      "1.0",
}

func initTracer(ctx context.Context, config OtelConfig) (func(context.Context) error, error) {

	tracerOpts := []sdktrace.TracerProviderOption{}

	exporterShutdownFunc := func(ctx context.Context) error { return nil }

	//exporter
	if config.ExporterOtlpEndpoint != "" {

		secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		if config.ExporterInsecure {
			secureOption = otlptracegrpc.WithInsecure()
		}

		exporter, err := otlptrace.New(
			ctx,
			otlptracegrpc.NewClient(
				secureOption,
				otlptracegrpc.WithEndpoint(config.ExporterOtlpEndpoint),
			),
		)
		if err != nil {
			return nil, err
		}

		tracerOpts = append(tracerOpts, sdktrace.WithBatcher(exporter))
		exporterShutdownFunc = exporter.Shutdown
	}

	//resources

	resources, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String(config.ServiceVersion),
			semconv.ServiceNamespaceKey.String(config.ServiceNamespace),
			semconv.DeploymentEnvironmentKey.String(config.DeploymentEnvironment),
			semconv.TelemetrySDKLanguageGo,
		),
	)
	if err != nil {
		return nil, err
	}

	tracerOpts = append(tracerOpts, sdktrace.WithResource(resources))

	//sampling
	sampler, err := buildTraceSampler(config.TracesSampler, config.TracesSamplerArg)
	if err != nil {
		return exporterShutdownFunc, err
	}
	tracerOpts = append(tracerOpts, sdktrace.WithSampler(sampler))

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(tracerOpts...),
	)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	return exporterShutdownFunc, nil
}

func Init(ctx context.Context, otelConfig OtelConfig) (func(context.Context) error, error) {

	err := otelConfig.sanitizeOtelConfig()
	if err != nil {
		return nil, err
	}

	cleanup, err := initTracer(ctx, otelConfig)
	if err != nil {
		return nil, err
	}

	return cleanup, nil
}

func AddTracingMiddlware(e *gin.Engine) {

	e.Use(otelgin.Middleware("alfred-server-name"))

}

func GetSpanContext(span trace.Span) (string, string, string) {

	traceId := span.SpanContext().TraceID().String()
	spanId := span.SpanContext().SpanID().String()
	flags := span.SpanContext().TraceFlags().String()

	return traceId, spanId, flags
}

func GetSpanFromContext(ctx context.Context) trace.Span {

	return trace.SpanFromContext(ctx)
}

func GetTracerFromContext(ctx context.Context) trace.Tracer {

	span := trace.SpanFromContext(ctx)
	return span.TracerProvider().Tracer(TracerName, trace.WithInstrumentationVersion(TracerVersion))

}

func SetSpanStatusError(span *trace.Span, err error) {

	(*span).RecordError(err)
	(*span).SetStatus(codes.Error, err.Error())

}

func (c *OtelConfig) sanitizeOtelConfig() error {

	if err := mergo.Merge(c, DEFAULT_OTEL_CONFIG); err != nil {
		return err
	}

	return nil
}

func buildTraceSampler(samplerStr string, samplerArgs string) (sdktrace.Sampler, error) {

	if strings.ToLower(samplerStr) == "always_on" {
		return sdktrace.AlwaysSample(), nil
	} else if strings.ToLower(samplerStr) == "always_off" {
		return sdktrace.NeverSample(), nil
	} else if strings.ToLower(samplerStr) == "parentbased_traceidratio" {

		ratio, err := strconv.ParseFloat(samplerArgs, 64)
		if err != nil {
			return sdktrace.AlwaysSample(), err
		}

		return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(ratio)), nil
	}

	return sdktrace.AlwaysSample(), errors.New("'" + samplerStr + "' not handled for trace sampling")
}
