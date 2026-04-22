package calculator

import (
	"pack-optimizer/internal/pack"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateOptimalPacks(t *testing.T) {
	p := pack.Pack{}
	calc := NewCalculator(p)

	tests := []struct {
		name       string
		totalItems int
		expected   []int
	}{
		{
			name:       "zero items, empty result",
			totalItems: 0,
			expected:   []int{},
		},
		{
			name:       "one item, smallest pack",
			totalItems: 1,
			expected:   []int{250},
		},
		{
			name:       "small quantity, smallest pack",
			totalItems: 250,
			expected:   []int{250},
		},
		{
			name:       "just over a pack, next size",
			totalItems: 251,
			expected:   []int{500},
		},
		{
			name:       "exact 500 pack",
			totalItems: 500,
			expected:   []int{500},
		},
		{
			name:       "750 items (500 + 250)",
			totalItems: 750,
			expected:   []int{500, 250},
		},
		{
			name:       "just under 1000",
			totalItems: 999,
			expected:   []int{1000},
		},
		{
			name:       "exact 1000 pack",
			totalItems: 1000,
			expected:   []int{1000},
		},
		{
			name:       "1250 items (1000 + 250)",
			totalItems: 1250,
			expected:   []int{1000, 250},
		},
		{
			name:       "medium quantity",
			totalItems: 501,
			expected:   []int{500, 250},
		},
		{
			name:       "exact 2000 pack",
			totalItems: 2000,
			expected:   []int{2000},
		},
		{
			name:       "2500 items (2000 + 500)",
			totalItems: 2500,
			expected:   []int{2000, 500},
		},
		{
			name:       "just under 5000",
			totalItems: 4999,
			expected:   []int{5000},
		},
		{
			name:       "exact 5000 pack",
			totalItems: 5000,
			expected:   []int{5000},
		},
		{
			name:       "just over 5000 (5001)",
			totalItems: 5001,
			expected:   []int{5000, 250},
		},
		{
			name:       "large quantity 3001",
			totalItems: 3001,
			expected:   []int{2000, 1000, 250},
		},
		{
			name:       "large quantity 3501",
			totalItems: 3501,
			expected:   []int{2000, 1000, 500, 250},
		},
		{
			name:       "double 5000 packs",
			totalItems: 10000,
			expected:   []int{5000, 5000},
		},
		{
			name:       "double 5000 plus 500",
			totalItems: 10500,
			expected:   []int{5000, 5000, 500},
		},
		{
			name:       "large quantity 12001",
			totalItems: 12001,
			expected:   []int{5000, 5000, 2000, 250},
		},
		{
			name:       "very large quantity 25000",
			totalItems: 25000,
			expected:   []int{5000, 5000, 5000, 5000, 5000},
		},
		{
			name:       "edge case 1500",
			totalItems: 1500,
			expected:   []int{1000, 500},
		},
		{
			name:       "edge case 3000",
			totalItems: 3000,
			expected:   []int{2000, 1000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculateOptimalPacks(tt.totalItems)
			assert.ElementsMatch(t, tt.expected, result, "CalculateOptimalPacks(%d) should return %v, but got %v", tt.totalItems, tt.expected, result)
		})
	}
}
