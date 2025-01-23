package mapper

import (
	"fmt"
	"goFlaky/core/strategy"
)

func CreateStrategy(id string) (strategy.IStrategy, error) {
	var strt strategy.IStrategy
	switch id {
	case "reverse":
		strt = strategy.ReverseOrder{}
	case "random":
		strt = strategy.RandomOrder{}
	default:
		return nil, fmt.Errorf("unsupported strategy: %s", id)
	}
	return strt, nil
}
