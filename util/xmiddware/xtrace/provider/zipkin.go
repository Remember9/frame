package provider

import (
	"context"
	"github.com/Remember9/frame/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func InitZipkinTracer(url string) func() {
	// Create Zipkin Exporter and install it as a global tracer.
	//
	// For demoing purposes, always sample. In a production application, you should
	// configure the sampler to a trace.ParentBased(trace.TraceIDRatioBased) set at the desired
	// ratio.
	name := "mysrv"
	if appConfig := config.GetAppConfig(); appConfig.Name != "" {
		name = appConfig.Name
		if appConfig.Version != "" {
			name = name + ":" + appConfig.Version
		}
		// carrier.Set(serviceHeader, appConfig.Name+"_"+appConfig.Version)
	}

	exporter, err := zipkin.New(
		url,
		// zipkin.WithLogger(logger),
		// zipkin.WithSDKOptions(sdktrace.WithSampler(sdktrace.AlwaysSample())),
	)
	if err != nil {
		// logger.Log(log.LevelError, "zipkin.New", err.Error())
		panic(err)
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(name),
			attribute.String("env", "dev"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return func() {
		_ = tp.Shutdown(context.Background())
	}
}
