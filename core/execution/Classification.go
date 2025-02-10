package execution

import (
	"goFlaky/core"
	"log"
)

type testId struct {
	Suite string
	Name  string
}

type classification struct {
	PreRunOutcome string
	OdRunOutcome  string
}

func CreateClassification(dj core.DependencyInjection, runId int, project string) {
	results, err := dj.Db.GetProjectTestResults(runId, project)
	if err != nil {
		log.Fatal(err)
	}

	testIds := make(map[testId]classification)
	for _, result := range results {
		resultTestId := testId{
			Suite: result.TestSuite,
			Name:  result.TestName,
		}
		resultClassification, ok := testIds[resultTestId]
		if !ok {
			resultClassification = classification{
				PreRunOutcome: "",
				OdRunOutcome:  "",
			}
		}
		if result.RunType == "PRE_RUN" {
			resultClassification.PreRunOutcome = testClassify(resultClassification.PreRunOutcome, result.TestOutcome)
		}
		if result.RunType == "OD_RUN" {
			resultClassification.OdRunOutcome = testClassify(resultClassification.OdRunOutcome, result.TestOutcome)
		}
		testIds[resultTestId] = resultClassification
	}

	for id, testClassification := range testIds {
		err := dj.Db.CreateClassification(
			runId,
			project,
			id.Suite,
			id.Name,
			generalClassify(testClassification.PreRunOutcome, testClassification.OdRunOutcome),
		)
		if err != nil {
			continue
		}
	}
}

func testClassify(prev string, outcome string) string {
	if prev == "" {
		return outcome
	}
	if prev == "FLAKY" || prev != outcome {
		return "FLAKY"
	}
	return prev
}

func generalClassify(preClassification string, odClassification string) string {
	if preClassification == "FLAKY" {
		return "PRE_RUN_FLAKY"
	}
	if odClassification == "FLAKY" {
		return "OD_RUN_FLAKY"
	}
	if preClassification != odClassification { //example: all pre runs fail but all od runs pass
		return "OD_RUN_FLAKY"
	}
	return preClassification
}
