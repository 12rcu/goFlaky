package framework

import "regexp"

type Config interface {
	Language() string
	TestAnnotation() *regexp.Regexp
	TestSuiteStart() *regexp.Regexp
	ImportStatement() *regexp.Regexp
	Imports() string
	IgnoreAnnotations() string
	ClassOrderAnnotation() string
	TestOrderAnnotations(index int) string
	IsTestContentForTestSuite(fileContent string, testSuite string, testName string) bool
}
