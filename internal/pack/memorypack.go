package pack

import (
	"context"
	"slices"
	"sync"
)

const maxPackSize = 10000

type InMemomorySvc struct {
	sync.Mutex
	sizes []int
}

func NewInMemorySvc(sizes []int) PackSvc {
	validate(sizes)
	return &InMemomorySvc{
		sizes: sort(sizes),
	}
}

func (p *InMemomorySvc) GetSizes(ctx context.Context) ([]int, error) {
	p.Lock()
	defer p.Unlock()
	return append([]int{}, p.sizes...), nil
}

func (p *InMemomorySvc) UpdateSizes(ctx context.Context, newSizes []int) error {
	if err := validate(newSizes); err != nil {
		return err
	}

	p.Lock()
	defer p.Unlock()
	p.sizes = sort(newSizes)
	return nil
}

func validate(sizes []int) error {
	if len(sizes) == 0 {
		return &ValidationError{"sizes cannot be empty"}
	}
	for _, size := range sizes {
		if size <= 0 {
			return &ValidationError{"size must be positive integer"}
		}
		if size > maxPackSize {
			return &ValidationError{"size must not exceed 10000"}
		}
	}
	return nil
}

func sort(sizes []int) []int {
	copy := append([]int{}, sizes...)
	slices.Sort(copy)
	return copy
}
