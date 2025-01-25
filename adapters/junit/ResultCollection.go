package junit

import (
	"encoding/xml"
	"goFlaky/adapters/util"
	"goFlaky/core"
	"goFlaky/core/framework"
	"os"
	"regexp"
)

func ResultCollection(resultPath string, dj core.DependencyInjection) []framework.TestResult {
	var testResults []framework.TestResult
	resRegex, err := regexp.Compile(`TEST-.+\.xml`)
	if err != nil {
		dj.TerminalLogChannel <- "[ERROR] while result collection " + err.Error()
		return testResults
	}
	files := util.SearchFile(resultPath, func(path string, fileName string) bool {
		return resRegex.MatchString(fileName)
	}, func(log string) {
		dj.TerminalLogChannel <- "[ERROR] while result collection " + log
	})
	for _, file := range files {
		xmlFile, err := os.Open(file)
		if err != nil {
			dj.TerminalLogChannel <- "[ERROR] while result collection " + err.Error()
			continue
		}

		byteValue, _ := os.ReadFile(file)
		var testSuite TestSuite
		err = xml.Unmarshal(byteValue, &testSuite)
		if err != nil {
			dj.TerminalLogChannel <- "[ERROR] while result collection " + err.Error()
			continue
		}

		for _, test := range testSuite.TestCases {
			outcome := "PASSED"
			if test.Failure != nil {
				outcome = "FAILURE"
			}
			testResults = append(testResults, framework.TestResult{
				TestSuite:   testSuite.Name,
				TestName:    test.Name,
				TestOutcome: outcome,
			})
		}

		err = xmlFile.Close()
		if err != nil {
			dj.TerminalLogChannel <- "[ERROR] while result collection " + err.Error()
		}
	}
	return testResults
}
