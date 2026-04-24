package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/aleksandarv/pack-optimizer/gen/optimizer"
	"github.com/aleksandarv/pack-optimizer/internal/calculator"
	"github.com/aleksandarv/pack-optimizer/internal/pack"
)

type optimizersrvc struct {
	psvc pack.PackSvc
	c    calculator.Calculator
}

func NewOptimizerSvc(p pack.PackSvc, c calculator.Calculator) optimizer.Service {
	return &optimizersrvc{
		psvc: p,
		c:    c,
	}
}

func (s *optimizersrvc) GetSizes(ctx context.Context) (res *optimizer.GetSizesResult, err error) {
	fmt.Printf("optimizer.getSizes")
	return &optimizer.GetSizesResult{
		Sizes: s.psvc.GetSizes(),
	}, nil
}

func (s *optimizersrvc) UpdateSizes(ctx context.Context, p *optimizer.UpdateSizesPayload) (err error) {
	log.Printf("optimizer.updateSizes")
	return s.psvc.UpdateSizes(p.Sizes)
}

func (s *optimizersrvc) Calculate(ctx context.Context, p *optimizer.CalculatePayload) (res *optimizer.CalculateResult, err error) {
	log.Printf("optimizer.calculate")
	var packs []*optimizer.Pack
	result := s.c.CalculateOptimalPacks(p.TotalItems)
	for _, pack := range result {
		packs = append(packs, &optimizer.Pack{
			Size:     pack.Size,
			Quantity: pack.Quantity,
		})
	}

	return &optimizer.CalculateResult{
		Packs: packs,
	}, nil
}
