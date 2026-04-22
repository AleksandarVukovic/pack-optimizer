package calculator

import (
	"fmt"
	"math"
	"pack-optimizer/internal/pack"
)

type Calculator interface {
	CalculateOptimalPacks(totalItems int) []int
}

type calculator struct {
	pack pack.Pack
}

func NewCalculator(p pack.Pack) Calculator {
	return &calculator{
		pack: p,
	}
}

func (c *calculator) CalculateOptimalPacks(totalItems int) []int {
	sizes := c.pack.GetSizes()
	result := []int{}

	// TODO: replace with logger
	fmt.Println()
	fmt.Printf("-------Calculating optimal packs for %d items with available sizes: %v\n", totalItems, sizes)

	if totalItems <= 0 {
		return result
	}

	r := c.optimizePacks(totalItems, sizes)
	fmt.Printf("Initial optimal packs: %v\n", r)
	// TODO: optimize generated packages
	result = c.optimizePacks(sum(r), sizes)
	fmt.Printf("Optimized packs: %v\n", result)
	return result
}

func (c *calculator) optimizePacks(totalItems int, sizes []int) []int {
	mresult := make(map[int]int)

	for totalItems > 0 {
		optimalSize := optimal(totalItems, sizes)
		if optimalSize == 0 {
			optimalSize = closest(totalItems, sizes)
		}

		numPacks := 1
		if found, n := c.findOptimalPack(totalItems, optimalSize, sizes); found {
			sizes = remove(optimalSize, sizes)
			numPacks = n
		}
		mresult[optimalSize] += numPacks
		totalItems -= numPacks * optimalSize
	}

	fmt.Printf("Returning result: %v\n", mresult)
	// TODO temp solution until we refactor
	result := make([]int, 0)
	for size, count := range mresult {
		for i := 0; i < count; i++ {
			result = append(result, size)
		}
	}
	return result
}

func (c *calculator) findOptimalPack(totalItems, packSize int, allowedSizes []int) (bool, int) {
	numPacks := totalItems / packSize
	for numPacks > 0 {
		remainingItems := totalItems - numPacks*packSize
		// we found optimal pack size
		if remainingItems == 0 {
			return true, numPacks
		}

		remainingSizes := remove(packSize, allowedSizes)
		if optimalSize := closest(remainingItems, remainingSizes); optimalSize > 0 {
			if found, _ := c.findOptimalPack(remainingItems, optimalSize, remainingSizes); found {
				// just check if there is optimal pack
				return true, numPacks
			}
		}
		numPacks--
	}
	return false, 0
}

func remove(num int, items []int) []int {
	result := make([]int, 0, len(items)-1)
	for _, item := range items {
		if item != num {
			result = append(result, item)
		}
	}
	return result
}

func closest(number int, sizes []int) int {
	if len(sizes) == 0 {
		return 0
	}
	closestSize := sizes[0]
	for _, size := range sizes {
		if math.Abs(float64(number-size)) < math.Abs(float64(number-closestSize)) {
			closestSize = size
		}
	}
	return closestSize
}

func sum(items []int) int {
	var sum int
	for _, item := range items {
		sum += item
	}
	return sum
}

func optimal(totalItems int, sizes []int) int {
	for _, size := range sizes {
		c := totalItems / size
		if c == 1 {
			return size
		}
	}
	return 0
}
