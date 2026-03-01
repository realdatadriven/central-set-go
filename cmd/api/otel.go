// https://opentelemetry.io/docs/languages/go/getting-started/

package main

/*import (
	"context"
	"errors"
	"time"

	"github.com/realdatadriven/central-set-go/internal/env"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error
	var err error
	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}
	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}
	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)
	// Set up trace provider.
	tracerProvider, err := newTracerProvider()
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)
	// Set up meter provider.
	meterProvider, err := newMeterProvider()
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)
	// Set up logger provider.
	loggerProvider, err := newLoggerProvider()
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)
	return shutdown, err
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider() (*trace.TracerProvider, error) {
	if env.GetString("OTEL_EXPORTER_OTLP_ENDPOINT", "") != "" {
		ctx := context.Background()
		exp, err := otlptracehttp.New(ctx)
		if err != nil {
			return nil, err
		}
		tracerProvider := trace.NewTracerProvider(trace.WithBatcher(exp))
		defer func() { tracerProvider.Shutdown(ctx) }()
		return tracerProvider, nil
	} else {
		traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
		tracerProvider := trace.NewTracerProvider(
			// trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(0.10))),
			trace.WithBatcher(traceExporter,
				// Default is 5s. Set to 1s for demonstrative purposes.
				trace.WithBatchTimeout(time.Second)),
		)
		return tracerProvider, nil
	}
}

func newMeterProvider() (*metric.MeterProvider, error) {
	if env.GetString("OTEL_EXPORTER_OTLP_ENDPOINT", "") != "" {
		ctx := context.Background()
		exp, err := otlpmetrichttp.New(ctx)
		if err != nil {
			return nil, err
		}
		meterProvider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(exp)))
		defer func() { meterProvider.Shutdown(ctx) }()
		return meterProvider, nil
	} else {
		metricExporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
		meterProvider := metric.NewMeterProvider(
			metric.WithReader(metric.NewPeriodicReader(metricExporter,
				// Default is 1m. Set to 3s for demonstrative purposes.
				metric.WithInterval(3*time.Second))),
		)
		return meterProvider, nil
	}
}

func newLoggerProvider() (*log.LoggerProvider, error) {
	if env.GetString("OTEL_EXPORTER_OTLP_ENDPOINT", "") != "" {
		ctx := context.Background()
		exp, err := otlploghttp.New(ctx)
		if err != nil {
			return nil, err
		}
		processor := log.NewBatchProcessor(exp)
		loggerProvider := log.NewLoggerProvider(log.WithProcessor(processor))
		defer func() { loggerProvider.Shutdown(ctx) }()
		return loggerProvider, nil
	} else {
		logExporter, err := stdoutlog.New(stdoutlog.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
		loggerProvider := log.NewLoggerProvider(
			log.WithProcessor(log.NewBatchProcessor(logExporter)),
		)
		return loggerProvider, nil
	}
}*/

/*
consider opentelemetry script (https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp):

package main

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
)

func main() {
	ctx := context.Background()
	exp, err := otlploghttp.New(ctx)
	if err != nil {
		panic(err)
	}

	processor := log.NewBatchProcessor(exp)
	provider := log.NewLoggerProvider(log.WithProcessor(processor))
	defer func() {
		if err := provider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	global.SetLoggerProvider(provider)

	// From here, the provider can be used by instrumentation to collect
	// telemetry.
}
after thhis how do i can i use the provider in the instrumentation to collect telemetry data
*/
