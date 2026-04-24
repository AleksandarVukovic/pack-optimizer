package pack

import (
	"errors"
	"slices"
	"sync"
)

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

func (p *InMemomorySvc) GetSizes() []int {
	p.Lock()
	defer p.Unlock()
	return append([]int{}, p.sizes...)
}

func (p *InMemomorySvc) UpdateSizes(newSizes []int) error {
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
		return errors.New("sizes cannot be empty")
	}
	for _, size := range sizes {
		if size <= 0 {
			return errors.New("size must be positive integer")
		}
	}
	return nil
}

func sort(sizes []int) []int {
	copy := append([]int{}, sizes...)
	slices.Sort(copy)
	return copy
}
