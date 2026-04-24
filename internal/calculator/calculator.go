package calculator

import (
	"fmt"
	"math"

	"github.com/aleksandarv/pack-optimizer/internal/pack"
)

type Calculator interface {
	CalculateOptimalPacks(totalItems int) []pack.Pack
}

type calculator struct {
	psvc pack.PackSvc
}

func NewCalculator(p pack.PackSvc) Calculator {
	return &calculator{
		psvc: p,
	}
}

func (c *calculator) CalculateOptimalPacks(totalItems int) []pack.Pack {
	sizes := c.psvc.GetSizes()
	result := make([]pack.Pack, 0)

	// TODO: replace with logger
	fmt.Println()
	fmt.Printf("-------Calculating optimal packs for %d items with available sizes: %v\n", totalItems, sizes)

	if totalItems <= 0 {
		return result
	}

	r := c.optimizePacks(totalItems, sizes)
	fmt.Printf("Initial optimal packs: %v\n", r)
	// TODO: optimize generated packages
	return c.optimizePacks(sum(r), sizes)
}

func (c *calculator) optimizePacks(totalItems int, sizes []int) []pack.Pack {
	result := make([]pack.Pack, 0)
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

	// TODO: replace with logger.debug
	fmt.Printf("Returning result: %v\n", mresult)
	for size, quantity := range mresult {
		result = append(result, pack.Pack{
			Size:     size,
			Quantity: quantity,
		})
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

func sum(packs []pack.Pack) int {
	sum := 0
	for _, pack := range packs {
		sum += pack.Size * pack.Quantity
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
