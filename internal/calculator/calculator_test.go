package calculator

import (
	"context"
	"testing"

	"github.com/aleksandarv/pack-optimizer/internal/logger"
	"github.com/aleksandarv/pack-optimizer/internal/pack"
)

func TestCalculateOptimalPacks(t *testing.T) {
	ctx := logger.WithCtx(context.Background(), logger.NewLogger(false))
	p := pack.NewInMemorySvc(pack.DefaultSizes)
	calc := NewCalculator(p)

	tests := []struct {
		name       string
		totalItems int
		expected   map[int]int
	}{
		{
			name:       "zero items, empty result",
			totalItems: 0,
			expected:   map[int]int{},
		},
		{
			name:       "one item, smallest pack",
			totalItems: 1,
			expected:   map[int]int{250: 1},
		},
		{
			name:       "small quantity, smallest pack",
			totalItems: 250,
			expected:   map[int]int{250: 1},
		},
		{
			name:       "just over a pack, next size",
			totalItems: 251,
			expected:   map[int]int{500: 1},
		},
		{
			name:       "exact 500 pack",
			totalItems: 500,
			expected:   map[int]int{500: 1},
		},
		{
			name:       "750 items (500 + 250)",
			totalItems: 750,
			expected:   map[int]int{500: 1, 250: 1},
		},
		{
			name:       "just under 1000",
			totalItems: 999,
			expected:   map[int]int{1000: 1},
		},
		{
			name:       "exact 1000 pack",
			totalItems: 1000,
			expected:   map[int]int{1000: 1},
		},
		{
			name:       "1250 items (1000 + 250)",
			totalItems: 1250,
			expected:   map[int]int{1000: 1, 250: 1},
		},
		{
			name:       "medium quantity",
			totalItems: 501,
			expected:   map[int]int{500: 1, 250: 1},
		},
		{
			name:       "exact 2000 pack",
			totalItems: 2000,
			expected:   map[int]int{2000: 1},
		},
		{
			name:       "2500 items (2000 + 500)",
			totalItems: 2500,
			expected:   map[int]int{2000: 1, 500: 1},
		},
		{
			name:       "just under 5000",
			totalItems: 4999,
			expected:   map[int]int{5000: 1},
		},
		{
			name:       "exact 5000 pack",
			totalItems: 5000,
			expected:   map[int]int{5000: 1},
		},
		{
			name:       "just over 5000 (5001)",
			totalItems: 5001,
			expected:   map[int]int{5000: 1, 250: 1},
		},
		{
			name:       "large quantity 3001",
			totalItems: 3001,
			expected:   map[int]int{2000: 1, 1000: 1, 250: 1},
		},
		{
			name:       "large quantity 3501",
			totalItems: 3501,
			expected:   map[int]int{2000: 1, 1000: 1, 500: 1, 250: 1},
		},
		{
			name:       "double 5000 packs",
			totalItems: 10000,
			expected:   map[int]int{5000: 2},
		},
		{
			name:       "double 5000 plus 500",
			totalItems: 10500,
			expected:   map[int]int{5000: 2, 500: 1},
		},
		{
			name:       "large quantity 12001",
			totalItems: 12001,
			expected:   map[int]int{5000: 2, 2000: 1, 250: 1},
		},
		{
			name:       "very large quantity 25000",
			totalItems: 25000,
			expected:   map[int]int{5000: 5},
		},
		{
			name:       "edge case 1500",
			totalItems: 1500,
			expected:   map[int]int{1000: 1, 500: 1},
		},
		{
			name:       "edge case 3000",
			totalItems: 3000,
			expected:   map[int]int{2000: 1, 1000: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculateOptimalPacks(ctx, tt.totalItems)
			if len(result) != len(tt.expected) {
				t.Errorf("CalculateOptimalPacks(%d) returned %d pack sizes, expected %d", tt.totalItems, len(result), len(tt.expected))
			}
			for _, pack := range result {
				expectedQuantity, exists := tt.expected[pack.Size]
				if !exists {
					t.Errorf("CalculateOptimalPacks(%d) returned unexpected pack size %d with count %d", tt.totalItems, pack.Size, pack.Quantity)
				} else if pack.Quantity != expectedQuantity {
					t.Errorf("CalculateOptimalPacks(%d) returned quantity %d for pack size %d, expected %d", tt.totalItems, pack.Quantity, pack.Size, expectedQuantity)
				}
			}
		})
	}
}

func TestCalculateOptimalPacks_WithHugeNumbers(t *testing.T) {
	ctx := logger.WithCtx(context.Background(), logger.NewLogger(false))

	const totalItems = 500000
	p := pack.NewInMemorySvc([]int{23, 31, 53})
	calc := NewCalculator(p)

	result := calc.CalculateOptimalPacks(ctx, totalItems)

	expected := map[int]int{
		23: 2,
		31: 7,
		53: 9429,
	}
	if len(result) != len(expected) {
		t.Errorf("CalculateOptimalPacks(%d) returned %d pack sizes, expected %d", totalItems, len(result), len(expected))
	}
	for _, pack := range result {
		expectedQuantity, exists := expected[pack.Size]
		if !exists {
			t.Errorf("CalculateOptimalPacks(%d) returned unexpected pack size %d with count %d", totalItems, pack.Size, pack.Quantity)
		} else if pack.Quantity != expectedQuantity {
			t.Errorf("CalculateOptimalPacks(%d) returned quantity %d for pack size %d, expected %d", totalItems, pack.Quantity, pack.Size, expectedQuantity)
		}
	}
}
