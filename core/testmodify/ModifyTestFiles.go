package testmodify

import (
	"goFlaky/core/framework"
)

func ModifyTestFiles(runOrders [][]int, testFileContent string, frConfig framework.Config) []string {
	var modifiedContentFiles []string
	for _, order := range runOrders {
		newContent := testFileContent
		newContent = AddFirstAnnotationBefore(newContent,
			frConfig.TestSuiteStart(),
			frConfig.ClassOrderAnnotation(),
		)
		newContent = AddAnnotationBefore(newContent, frConfig.TestAnnotation(), func(matchNum int) string {
			return frConfig.TestOrderAnnotations(order[matchNum])
		})
		newContent = AddImports(newContent, frConfig.Imports(), frConfig.ImportStatement())
		modifiedContentFiles = append(modifiedContentFiles, newContent)
	}
	return modifiedContentFiles
}
