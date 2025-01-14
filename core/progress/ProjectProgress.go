package progress

import (
	"errors"
	"slices"
)

type ProjectProgress struct {
	Identifier string
	Status     string
	Runs       int
	Index      int
}

func GetProgressIndex(identifier string, progress []ProjectProgress) (int, error) {
	progressIndex := slices.IndexFunc(progress, func(e ProjectProgress) bool {
		return e.Identifier == identifier
	})

	if progressIndex == -1 {
		return 0, errors.New("project not found")
	}
	return progressIndex, nil
}

func ProgressProject(identifier string, progress []ProjectProgress) error {
	progressIndex, err := GetProgressIndex(identifier, progress)

	if err != nil {
		return err
	}

	progress[progressIndex].Index++
	return nil
}
