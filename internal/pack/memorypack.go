package pack

import (
	"errors"
	"slices"
	"sync"
)

type InMemomoryPack struct {
	sync.Mutex
	sizes []int
}

func NewInMemoryPack(sizes []int) Pack {
	validate(sizes)
	return &InMemomoryPack{
		sizes: sort(sizes),
	}
}

func (p *InMemomoryPack) GetSizes() []int {
	p.Lock()
	defer p.Unlock()
	return append([]int{}, p.sizes...)
}

func (p *InMemomoryPack) UpdateSizes(newSizes []int) error {
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
