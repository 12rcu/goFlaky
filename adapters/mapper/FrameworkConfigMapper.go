package mapper

import (
	"fmt"
	"goFlaky/adapters/junit"
	"goFlaky/core/framework"
)

func CreateFrameworkConfig(id string) (framework.Config, error) {
	var config framework.Config
	switch id {
	case "junit":
		config = junit.CreateNew()
	default:
		return nil, fmt.Errorf("unsupported framework: %s", id)
	}
	return config, nil
}
