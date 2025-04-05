package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/hollowdll/go-grpc-microservices/services/inventory/config"
	"github.com/hollowdll/go-grpc-microservices/services/inventory/internal/adapters/db"
	"github.com/hollowdll/go-grpc-microservices/services/inventory/internal/adapters/grpc"
	"github.com/hollowdll/go-grpc-microservices/services/inventory/internal/application/core/api"
)

func initApplication(application *api.Application, cfg *config.Config) {
	if cfg.IsDevelopmentMode() {
		log.Println("development mode detected: populating test data ...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := application.PopulateTestData(ctx); err != nil {
			log.Fatalf("failed to populate test data: %v", err)
		}
	}
}

func main() {
	log.Println("starting inventory service ...")

	config.InitConfig()
	cfg := config.LoadConfig()
	log.Printf("running application in %s mode", cfg.ApplicationMode)

	// Tracing
	ctx := context.Background()
	traceExporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint("otel-collector:4318"), otlptracehttp.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("inventory-service"),
		)),
	)
	otel.SetTracerProvider(tp)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shut down tracer provider: %v", err)
		}
	}()

	dbAdapter, err := db.NewPostgresAdapter(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	application := api.NewApplication(dbAdapter)
	initApplication(application, cfg)

	grpcAdapter := grpc.NewAdapter(application, cfg)
	grpcAdapter.Run()
}
