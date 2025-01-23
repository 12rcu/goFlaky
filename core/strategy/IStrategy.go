package strategy

type IStrategy interface {
	// GenerateOrder get the maximum number of tests per file and returns all orders the test suite should be executed
	GenerateOrder(numTests int) [][]int
}
