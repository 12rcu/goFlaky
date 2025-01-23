package strategy

import (
	"math/rand"
)

type RandomOrder struct{}

func (RandomOrder) GenerateOrder(numTests int) [][]int {
	var orders [][]int
	for i := 0; i < numTests; i++ {
		var temp []int
		for j := 0; j < numTests; j++ {
			temp = append(temp, j)
		}
		shuffleSlice(temp)
		orders = append(orders, temp)
	}
	return orders
}

func shuffleSlice(order []int) {
	// Fisher–Yates shuffle
	rand.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })
}
