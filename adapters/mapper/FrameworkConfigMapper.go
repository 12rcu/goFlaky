package mapper

import (
	"fmt"
	"goFlaky/adapters/junit"
	"goFlaky/core/framework"
	"strings"
)

func CreateFrameworkConfig(id string) (framework.Config, error) {
	var config framework.Config
	switch strings.ToLower(id) {
	case "junit":
		config = junit.CreateNew()
		break
	default:
		return nil, fmt.Errorf("unsupported framework: %s", id)
	}
	return config, nil
}
