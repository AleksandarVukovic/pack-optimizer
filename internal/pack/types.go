package pack

var (
	DefaultSizes = []int{250, 500, 1000, 2000, 5000}
)

type Pack interface {
	GetSizes() []int
	UpdateSizes(newSizes []int) error
}
