package pack

var (
	sizes = []int{250, 500, 1000, 2000, 5000}
)

type Pack struct {
}

// when we replace const array with dynamic config, don't forget to SORT items

func (p Pack) GetSizes() []int {
	return append([]int(nil), sizes...)
}
