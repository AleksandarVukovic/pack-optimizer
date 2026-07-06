package telemetry

import "context"

const (
	ServiceName       = "pack-optimizer"
	APIComponentName  = "pack-optimizer/api"
	PackComponentName = "pack-optimizer/pack"

	CalculateCountMetric             = "optimizer.calculate.count"
	UpdatePackSizesCountMetric       = "optimizer.update_pack_sizes.count"
	UpdatePackSizesFailedCountMetric = "optimizer.update_pack_sizes.failed.count"
)

type Attribute struct {
	Key   string
	Value string
}

type Tracer interface {
	Trace(ctx context.Context, method string, attrs ...Attribute) (context.Context, func(err *error))
}

type Meter interface {
	Increment(ctx context.Context, name string)
}
