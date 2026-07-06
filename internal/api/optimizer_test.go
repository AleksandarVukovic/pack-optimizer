package api

import (
	"context"
	"errors"
	"testing"

	goaoptimizer "github.com/aleksandarv/pack-optimizer/gen/optimizer"
	"github.com/aleksandarv/pack-optimizer/internal/logger"
	"github.com/aleksandarv/pack-optimizer/internal/pack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetPackSizes(t *testing.T) {
	tests := []struct {
		name          string
		sizes         []int
		getErr        error
		expectedSizes []int
		expectErr     bool
	}{
		{
			name:          "returns configured sizes",
			sizes:         []int{250, 500, 1000},
			expectedSizes: []int{250, 500, 1000},
		},
		{
			name:          "returns empty slice",
			sizes:         []int{},
			expectedSizes: []int{},
		},
		{
			name:      "returns error when service fails",
			getErr:    errors.New("db unavailable"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			psvc := &mockPackSvc{}
			psvc.On("GetSizes", mock.Anything).Return(tt.sizes, tt.getErr)

			svc := NewOptimizerSvc(psvc, &mockCalculator{}, newMockMeter())
			res, err := svc.GetPackSizes(newCtx(t))

			if tt.expectErr {
				require.Error(t, err)
				assert.Nil(t, res)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedSizes, res.Sizes)
			}
			psvc.AssertExpectations(t)
		})
	}
}

func TestUpdatePackSizes(t *testing.T) {
	tests := []struct {
		name      string
		sizes     []int
		updateErr error
		expectErr bool
	}{
		{
			name:      "updates sizes successfully",
			sizes:     []int{100, 200, 300},
			updateErr: nil,
			expectErr: false,
		},
		{
			name:      "returns error on empty sizes",
			sizes:     []int{},
			updateErr: &pack.ValidationError{Msg: "sizes cannot be empty"},
			expectErr: true,
		},
		{
			name:      "returns error on non-positive size",
			sizes:     []int{-100},
			updateErr: &pack.ValidationError{Msg: "size must be positive integer"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			psvc := &mockPackSvc{}
			psvc.On("UpdateSizes", mock.Anything, tt.sizes).Return(tt.updateErr)

			svc := NewOptimizerSvc(psvc, &mockCalculator{}, newMockMeter())
			err := svc.UpdatePackSizes(newCtx(t), &goaoptimizer.UpdatePackSizesPayload{
				Sizes: tt.sizes,
			})

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			psvc.AssertExpectations(t)
		})
	}
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name          string
		quantity      int
		calcResult    []pack.Pack
		expectedPacks map[int]int
	}{
		{
			name:     "returns correct packs for large quantity",
			quantity: 12001,
			calcResult: []pack.Pack{
				{Size: 5000, Quantity: 2},
				{Size: 2000, Quantity: 1},
				{Size: 250, Quantity: 1},
			},
			expectedPacks: map[int]int{5000: 2, 2000: 1, 250: 1},
		},
		{
			name:     "returns correct packs for small quantity",
			quantity: 501,
			calcResult: []pack.Pack{
				{Size: 500, Quantity: 1},
				{Size: 250, Quantity: 1},
			},
			expectedPacks: map[int]int{500: 1, 250: 1},
		},
		{
			name:          "returns empty packs for zero quantity",
			quantity:      0,
			calcResult:    []pack.Pack{},
			expectedPacks: map[int]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newCtx(t)
			calc := &mockCalculator{}
			calc.On("CalculateOptimalPacks", ctx, tt.quantity).Return(tt.calcResult, nil)

			svc := NewOptimizerSvc(&mockPackSvc{}, calc, newMockMeter())
			res, err := svc.Calculate(ctx, &goaoptimizer.CalculatePayload{Quantity: tt.quantity})

			require.NoError(t, err)
			require.Len(t, res.Packs, len(tt.expectedPacks))

			packMap := make(map[int]int)
			for _, p := range res.Packs {
				packMap[p.Size] = p.Quantity
			}
			assert.Equal(t, tt.expectedPacks, packMap)
			calc.AssertExpectations(t)
		})
	}
}

func newCtx(t *testing.T) context.Context {
	t.Helper()
	return logger.WithCtx(context.Background(), logger.NewLogger(false))
}

type mockPackSvc struct {
	mock.Mock
}

func (m *mockPackSvc) GetSizes(ctx context.Context) ([]int, error) {
	args := m.Called(ctx)
	return args.Get(0).([]int), args.Error(1)
}

func (m *mockPackSvc) UpdateSizes(ctx context.Context, sizes []int) error {
	args := m.Called(ctx, sizes)
	return args.Error(0)
}

type mockCalculator struct {
	mock.Mock
}

func (m *mockCalculator) CalculateOptimalPacks(ctx context.Context, totalItems int) ([]pack.Pack, error) {
	args := m.Called(ctx, totalItems)
	return args.Get(0).([]pack.Pack), args.Error(1)
}

type mockMeter struct {
	mock.Mock
}

func (m *mockMeter) Increment(ctx context.Context, name string) {
	m.Called(ctx, name)
}

func newMockMeter() *mockMeter {
	m := &mockMeter{}
	m.On("Increment", mock.Anything, mock.Anything).Return()
	return m
}
