package testmodify

import (
	"goFlaky/core/framework"
)

func ModifyTestFile(runOrder []int, testFileContent string, frConfig framework.Config) string {
	newContent := testFileContent
	newContent = AddFirstAnnotationBefore(newContent,
		frConfig.TestSuiteStart(),
		frConfig.ClassOrderAnnotation(),
	)
	newContent = AddAnnotationBefore(newContent, frConfig.TestAnnotation(), func(matchNum int) string {
		return frConfig.TestOrderAnnotations(runOrder[matchNum])
	})
	newContent = AddImports(newContent, frConfig.Imports(), frConfig.ImportStatement())
	return newContent
}
