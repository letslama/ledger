package api

import (
	_ "embed"

	"github.com/formancehq/ledger/internal/api/backend"
	"github.com/formancehq/ledger/internal/engine"
	"github.com/formancehq/ledger/internal/opentelemetry/metrics"
	"github.com/formancehq/ledger/internal/storage/driver"
	"github.com/formancehq/stack/libs/go-libs/health"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/fx"
)

type Config struct {
	Version string
}

func Module(cfg Config) fx.Option {
	return fx.Options(
		fx.Provide(NewRouter),
		fx.Provide(func(storageDriver *driver.Driver, resolver *engine.Resolver) backend.Backend {
			return backend.NewDefaultBackend(storageDriver, cfg.Version, resolver)
		}),
		fx.Provide(fx.Annotate(noop.NewMeterProvider, fx.As(new(metric.MeterProvider)))),
		fx.Decorate(fx.Annotate(func(meterProvider metric.MeterProvider) (metrics.GlobalRegistry, error) {
			return metrics.RegisterGlobalRegistry(meterProvider)
		}, fx.As(new(metrics.GlobalRegistry)))),
		health.Module(),
	)
}
