package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log"
	nooplog "go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/metric"
	noopmetric "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
	nooptrace "go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// The value of the ServiceVersion attribute associated with the otel
// Resource. This should be overwritten at build time:
//
//	export SERVICE_VERSION=$(git rev-parse --short HEAD)
//	go build -o out \
//	  -X github.com/cmsgov/emmy-api/pkg/core.ServiceVersion=${SERVICE_VERSION}
var ServiceVersion = "UNSET"

// OtelService provides methods for handling OpenTelemetry operations.
type OtelService interface {
	// Return the current span associated with context c.
	SpanFromContext(c context.Context) trace.Span
	// Return the underlying LoggerProvider
	LoggerProvider() log.LoggerProvider
	// Cleanup any resources associated with this object. Errors
	// are logged with logger.Error or logger.ErrorContext and
	// then dismissed.
	Shutdown(c context.Context, logger *slog.Logger)
}

// Initialize the appropriate implementation of OtelService. If
// cfg.Disabled then the implementation uses the otel/*/noop packages.
func NewOtelService(ctx context.Context, cfg *Config) (OtelService, error) {
	if cfg.Otel.Disable {
		return newOtelServiceNoop(), nil
	}

	return newOtelServiceGRPC(ctx, cfg)
}

type shutdownFn = func(context.Context) error

type otelService struct {
	meterProvider  metric.MeterProvider
	tracerProvider trace.TracerProvider
	logProvider    log.LoggerProvider
	shutdown       shutdownFn
}

func (s otelService) LoggerProvider() log.LoggerProvider {
	return s.logProvider
}

func (otelService) SpanFromContext(c context.Context) trace.Span {
	return trace.SpanFromContext(c)
}

func (s otelService) Shutdown(c context.Context, logger *slog.Logger) {
	shutdownCtx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	err := s.shutdown(shutdownCtx)
	if err != nil {
		logger.ErrorContext(
			c,
			"Error shutting down otel",
			"err",
			err,
		)
	}
}

var _ OtelService = (*otelService)(nil)

var noopShutdown = func(_ context.Context) error { return nil }

func newOtelServiceNoop() OtelService {
	return otelService{
		meterProvider:  noopmetric.NewMeterProvider(),
		tracerProvider: nooptrace.NewTracerProvider(),
		logProvider:    nooplog.NewLoggerProvider(),
		shutdown:       noopShutdown,
	}
}

func newOtelServiceGRPC(ctx context.Context, cfg *Config) (OtelService, error) {
	res, err := newResource(ctx)
	if err != nil {
		return nil, err
	}

	conn, err := newConn(cfg)
	if err != nil {
		return nil, err
	}

	traceProvider, shutdownTraceProvider, err := newTraceProvider(ctx, res, conn)
	if err != nil {
		return nil, err
	}

	meterProvider, shutdownMeterProvider, err := newMeterProvider(ctx, res, conn)
	if err != nil {
		return nil, err
	}

	logProvider := sdklog.NewLoggerProvider(sdklog.WithResource(res))

	shutdown := func(shutdownCtx context.Context) error {
		return errors.Join(
			shutdownTraceProvider(shutdownCtx),
			shutdownMeterProvider(shutdownCtx),
		)
	}

	return otelService{
		meterProvider:  meterProvider,
		tracerProvider: traceProvider,
		logProvider:    logProvider,
		shutdown:       shutdown,
	}, nil
}

func newConn(cfg *Config) (*grpc.ClientConn, error) {
	conn, e := grpc.NewClient(
		cfg.Otel.OtlpExporter.Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if e != nil {
		return nil, fmt.Errorf("failed to create GRPC connection: %w", e)
	}

	return conn, nil
}

func newResource(ctx context.Context) (*resource.Resource, error) {
	res, e := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceVersion(ServiceVersion),
		),
	)

	if e == nil {
		return res, nil
	}
	return nil, fmt.Errorf("failed to create telemetry resource: %w", e)
}

func newTraceProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (trace.TracerProvider, shutdownFn, error) {
	exp, e := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithGRPCConn(conn),
	)

	if e != nil {
		return nil, nil, fmt.Errorf("failed to create trace exporter: %w", e)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return provider, provider.Shutdown, nil
}

func newMeterProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (metric.MeterProvider, shutdownFn, error) {
	exp, e := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithGRPCConn(conn),
	)

	if e != nil {
		return nil, nil, fmt.Errorf("failed to create meter exporter: %w", e)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp)),
		sdkmetric.WithResource(res),
	)

	return provider, provider.Shutdown, nil
}
