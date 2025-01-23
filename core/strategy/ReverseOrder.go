package strategy

type ReverseOrder struct{}

func (ReverseOrder) GenerateOrder(numTests int) [][]int {
	var normal []int
	var reverse []int

	for i := 1; i <= numTests; i++ {
		normal = append(normal, i)
		reverse = append(reverse, numTests-(i-1))
	}

	return [][]int{normal, reverse}
}
