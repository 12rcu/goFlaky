package junit

import (
	"encoding/xml"
	"goFlaky/adapters/util"
	"goFlaky/core/framework"
	"log"
	"os"
	"regexp"
)

func ResultCollection(resultPath string) []framework.TestResult {
	var testResults []framework.TestResult
	resRegex, err := regexp.Compile(`TEST-.+\.xml`)
	if err != nil {
		log.Default().Println(err)
		return testResults
	}
	files := util.SearchFile(resultPath, func(path string, fileName string) bool {
		return resRegex.MatchString(fileName)
	})
	for _, file := range files {
		xmlFile, err := os.Open(file)
		if err != nil {
			log.Default().Println(err)
			continue
		}

		byteValue, _ := os.ReadFile(file)
		var testSuite TestSuite
		err = xml.Unmarshal(byteValue, &testSuite)
		if err != nil {
			log.Default().Println(err)
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
			log.Default().Println(err)
		}
	}
	return testResults
}
