package pack

import "context"

var (
	DefaultSizes = []int{250, 500, 1000, 2000, 5000}
)

type PackSvc interface {
	GetSizes(ctx context.Context) ([]int, error)
	UpdateSizes(ctx context.Context, newSizes []int) error
}

type Pack struct {
	Size     int
	Quantity int
}

type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string {
	return e.Msg
}
