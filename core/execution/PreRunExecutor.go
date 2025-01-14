package execution

import "goFlaky/core"

func PreRunExecute(project core.ConfigProject, dj core.DependencyInjection) {

	for i := 0; i < int(project.PreRuns); i++ {

	}
}
