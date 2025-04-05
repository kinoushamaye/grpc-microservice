package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/hollowdll/go-grpc-microservices/services/payment/config"
	"github.com/hollowdll/go-grpc-microservices/services/payment/internal/adapters/grpc"
	"github.com/hollowdll/go-grpc-microservices/services/payment/internal/application/core/api"
)

func main() {
	log.Println("starting payment service ...")

	cfg := config.NewConfig()
	log.Printf("application is running in %s mode", cfg.ApplicationMode)

	ctx := context.Background()

	// Tracing
	traceExporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint("otel-collector:4318"), otlptracehttp.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("payment-service"),
		)),
	)
	otel.SetTracerProvider(tp)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shut down tracer provider: %v", err)
		}
	}()

	application := api.NewApplication()
	grpcAdapter := grpc.NewAdapter(application, cfg)
	grpcAdapter.Run()
}
