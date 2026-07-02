package api

import (
	"context"
	"errors"

	goaoptimizer "github.com/aleksandarv/pack-optimizer/gen/optimizer"
	"github.com/aleksandarv/pack-optimizer/internal/calculator"
	"github.com/aleksandarv/pack-optimizer/internal/logger"
	"github.com/aleksandarv/pack-optimizer/internal/pack"
	"github.com/aleksandarv/pack-optimizer/internal/telemetry"
)

type optimizersrvc struct {
	psvc pack.PackSvc
	c    calculator.Calculator
	m    telemetry.Meter
}

func NewOptimizerSvc(p pack.PackSvc, c calculator.Calculator, m telemetry.Meter) goaoptimizer.Service {
	return &optimizersrvc{
		psvc: p,
		c:    c,
		m:    m,
	}
}

func (s *optimizersrvc) Health(context.Context) (res *goaoptimizer.HealthResult, err error) {
	return &goaoptimizer.HealthResult{
		Status: "ok",
	}, nil
}

func (s *optimizersrvc) GetPackSizes(ctx context.Context) (res *goaoptimizer.GetPackSizesResult, err error) {
	log := logger.FromCtx(ctx)
	log.Info("optimizer.getSizes")

	return &goaoptimizer.GetPackSizesResult{
		Sizes: s.psvc.GetSizes(),
	}, nil
}

func (s *optimizersrvc) UpdatePackSizes(ctx context.Context, p *goaoptimizer.UpdatePackSizesPayload) (err error) {
	log := logger.FromCtx(ctx)
	log.Info("optimizer.updateSizes")

	s.m.Increment(ctx, telemetry.UpdatePackSizesCountMetric)

	if err = s.psvc.UpdateSizes(p.Sizes); err != nil {
		log.Error("error while updating pack sizes", "error", err.Error())
		s.m.Increment(ctx, telemetry.UpdatePackSizesFailedCountMetric)
		if verr, ok := errors.AsType[*pack.ValidationError](err); ok {
			return goaoptimizer.MakeBadRequest(verr)
		}
		return goaoptimizer.MakeInternalServerError(err)
	}
	return nil
}

func (s *optimizersrvc) Calculate(ctx context.Context, p *goaoptimizer.CalculatePayload) (res *goaoptimizer.CalculateResult, err error) {
	log := logger.FromCtx(ctx)
	log.Info("optimizer.calculate")

	s.m.Increment(ctx, telemetry.CalculateCountMetric)

	var packs []*goaoptimizer.Pack
	result := s.c.CalculateOptimalPacks(ctx, p.Quantity)
	for _, pack := range result {
		packs = append(packs, &goaoptimizer.Pack{
			Size:     pack.Size,
			Quantity: pack.Quantity,
		})
	}

	return &goaoptimizer.CalculateResult{
		Packs: packs,
	}, nil
}
