package junit

import (
	"goFlaky/core/framework"
	"regexp"
	"strconv"
	"strings"
)

type JUnit string

func CreateNew() framework.Config {
	return JUnit("kotlin")
}

func (j JUnit) Language() string {
	return string(j)
}

func (j JUnit) TestAnnotation() *regexp.Regexp {
	return regexp.MustCompile(`@Test(\s|\n|\r\n)`)
}

func (j JUnit) TestSuiteStart() *regexp.Regexp {
	return regexp.MustCompile(`(public |private |protected |internal )?class`)
}

func (j JUnit) ImportStatement() *regexp.Regexp {
	return regexp.MustCompile(`package (.*)`)
}

func (j JUnit) Imports() string {
	if string(j) == "kotlin" {
		return "import org.junit.jupiter.api.Disabled\n" +
			"import org.junit.jupiter.api.Order\n" +
			"import org.junit.jupiter.api.MethodOrderer\n" +
			"import org.junit.jupiter.api.TestMethodOrder\n"

	} else {
		return "import org.junit.jupiter.api.Disabled\n" +
			"import org.junit.jupiter.api.Order\n" +
			"import org.junit.jupiter.api.MethodOrderer\n" +
			"import org.junit.jupiter.api.TestMethodOrder\n"
	}
}

func (j JUnit) ClassOrderAnnotation() string {
	if string(j) == "kotlin" {
		return "@TestMethodOrder(MethodOrderer.OrderAnnotation::class)"
	} else {
		return "@TestMethodOrder(MethodOrderer.OrderAnnotation.class)"
	}
}

func (j JUnit) IgnoreAnnotations() string {
	return "@Disabled(\"kFlaky ignore\")"
}

func (j JUnit) TestOrderAnnotations(index int) string {
	return "@Order(" + strconv.Itoa(index) + ")"
}

func (j JUnit) IsTestContentForTestSuite(fileContent string, testSuite string, _ string) bool {
	// Split the test suite string by "."
	testSuiteL := strings.Split(testSuite, ".")

	// Get the package ID and class name
	var packageId, className string
	if len(testSuiteL) > 1 {
		packageId = strings.Join(testSuiteL[:len(testSuiteL)-1], "\\.")
		className = testSuiteL[len(testSuiteL)-1]
	} else {
		className = testSuiteL[0]
	}

	// Check if the class is present in the test file content
	clazz := containsClass(fileContent, className)
	pkg := true
	if packageId != "" {
		pkg = containsPackage(fileContent, packageId)
	}

	// Return true if both class and package are found
	return clazz && pkg
}

// Helper function to check if class name is present in the test file content
func containsClass(testFileContent string, className string) bool {
	rgx := regexp.MustCompile(`(public |private |protected |internal )?class ` + className)
	return rgx.MatchString(testFileContent)
}

// Helper function to check if package name is present in the test file content
func containsPackage(testFileContent string, packageId string) bool {
	rgx := regexp.MustCompile(`package ` + packageId)
	return rgx.MatchString(testFileContent)
}
