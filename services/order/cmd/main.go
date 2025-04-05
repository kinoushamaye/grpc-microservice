package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/hollowdll/go-grpc-microservices/services/order/config"
	"github.com/hollowdll/go-grpc-microservices/services/order/internal/adapters/grpc"
	"github.com/hollowdll/go-grpc-microservices/services/order/internal/adapters/inventory"
	"github.com/hollowdll/go-grpc-microservices/services/order/internal/adapters/payment"
	"github.com/hollowdll/go-grpc-microservices/services/order/internal/application/core/api"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initTelemetry() func() {
	exporter, err := otlptracehttp.New(context.Background(), otlptracehttp.WithInsecure(), otlptracehttp.WithEndpoint("otel-collector:4318"))
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("order-service"),
		)),
	)
	otel.SetTracerProvider(tp)

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("error shutting down tracer provider: %v", err)
		}
	}
}

func main() {
	log.Println("starting order service ...")

	shutdown := initTelemetry()
	defer shutdown()

	config.InitConfig()
	cfg := config.NewConfig()
	log.Printf("running application in %s mode", cfg.ApplicationMode)

	inventoryAdapter, err := inventory.NewAdapter(cfg)
	if err != nil {
		log.Fatalf("failed to initialize gRPC client for inventory service: %v", err)
	}
	defer inventoryAdapter.CloseConnection()

	paymentAdapter, err := payment.NewAdapter(cfg)
	if err != nil {
		log.Fatalf("failed to initialize gRPC client for payment service: %v", err)
	}
	defer paymentAdapter.CloseConnection()

	application := api.NewApplication(inventoryAdapter, paymentAdapter)
	grpcAdapter := grpc.NewAdapter(application, cfg)
	grpcAdapter.Run()
}
