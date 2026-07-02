package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
)

type otelMeter struct {
	meterName string
}

func NewMeter(name string) Meter {
	return &otelMeter{meterName: name}
}

func (m *otelMeter) Increment(ctx context.Context, name string) {
	counter, err := otel.Meter(m.meterName).Int64Counter(name)
	if err != nil {
		return
	}
	counter.Add(ctx, 1)
}
