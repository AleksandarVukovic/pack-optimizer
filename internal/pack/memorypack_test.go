package pack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInMemorySvc(t *testing.T) {
	tests := []struct {
		name     string
		sizes    []int
		expected []int
	}{
		{
			name:     "empty slice",
			sizes:    []int{},
			expected: []int{},
		},
		{
			name:     "single size",
			sizes:    []int{250},
			expected: []int{250},
		},
		{
			name:     "multiple sizes sorted",
			sizes:    []int{250, 500, 1000},
			expected: []int{250, 500, 1000},
		},
		{
			name:     "multiple sizes unsorted",
			sizes:    []int{1000, 250, 500},
			expected: []int{250, 500, 1000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewInMemorySvc(tt.sizes)
			assert.NotNil(t, svc)
			assert.IsType(t, &InMemomorySvc{}, svc)
			result := svc.GetSizes()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewInMemoryPack_DefensiveCopy(t *testing.T) {
	original := []int{250, 500, 1000}
	svc := NewInMemorySvc(original)
	original[0] = 555
	result := svc.GetSizes()

	assert.Equal(t, []int{250, 500, 1000}, result)
}

func TestUpdateSizes(t *testing.T) {
	tests := []struct {
		name         string
		initialSizes []int
		newSizes     []int
		expected     []int
		expectError  bool
	}{
		{
			name:         "update from empty",
			initialSizes: []int{},
			newSizes:     []int{250, 500},
			expected:     []int{250, 500},
			expectError:  false,
		},
		{
			name:         "update to single size",
			initialSizes: []int{250, 500, 1000},
			newSizes:     []int{750},
			expected:     []int{750},
			expectError:  false,
		},
		{
			name:         "update to empty should fail",
			initialSizes: []int{250, 500},
			newSizes:     []int{},
			expected:     []int{250, 500},
			expectError:  true,
		},
		{
			name:         "update with different sizes",
			initialSizes: []int{100, 200},
			newSizes:     []int{300, 400, 500},
			expected:     []int{300, 400, 500},
			expectError:  false,
		},
		{
			name:         "update with unsorted sizes",
			initialSizes: []int{100, 200},
			newSizes:     []int{500, 300, 400},
			expected:     []int{300, 400, 500},
			expectError:  false,
		},
		{
			name:         "update with negative size should fail",
			initialSizes: []int{250, 500},
			newSizes:     []int{-100},
			expected:     []int{250, 500},
			expectError:  true,
		},
		{
			name:         "update with zero size should fail",
			initialSizes: []int{250, 500},
			newSizes:     []int{0},
			expected:     []int{250, 500},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pack := NewInMemorySvc(tt.initialSizes)
			err := pack.UpdateSizes(tt.newSizes)

			if tt.expectError {
				assert.Error(t, err)
				// verify that sizes did not change
				result := pack.GetSizes()
				assert.Equal(t, tt.expected, result)
			} else {
				assert.NoError(t, err)
				result := pack.GetSizes()
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestUpdateSizes_DefensiveCopy(t *testing.T) {
	svc := NewInMemorySvc([]int{250})
	newSizes := []int{500, 1000, 2000}

	err := svc.UpdateSizes(newSizes)
	assert.NoError(t, err)

	newSizes[0] = 9999
	result := svc.GetSizes()
	assert.Equal(t, []int{500, 1000, 2000}, result)
}

func TestInMemorySvcImplementsPackInterface(t *testing.T) {
	svc := NewInMemorySvc([]int{250, 500})
	var _ PackSvc = svc
	assert.NotNil(t, svc)
}
