package util

import (
	"os"
	"path/filepath"
	"regexp"
)

func SearchFile(path string, expr func(path string, fileName string) bool, logErr func(log string)) []string {
	files := make([]string, 0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if expr(path, info.Name()) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		logErr(err.Error())
	}
	return files
}

// SearchFileWithContent Searches a given path for files that matches a given expression.
//
//	SearchFileWithContent(
//		"/home/exm/test",
//		func(path string, fileName string, fileContent string) bool {
//			return strings.Contains(fileContent, "myAwesomeString")
//		},
//	)
func SearchFileWithContent(path string, expr func(path string, fileName string, fileContent string) bool, logErr func(log string)) []string {
	files := make([]string, 0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if expr(path, info.Name(), string(content)) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		logErr(err.Error())
	}
	return files
}

// SearchFilesAndCount searches all files in a given directory and counts how often the regex is matched in a file.
// A list with the count of matches in each file is returned alongside with the highest number of matches.
func SearchFilesAndCount(path string, expr *regexp.Regexp, logErr func(log string)) ([]int, int) {
	counts := make([]int, 0)
	maxMatches := 0
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		matches := len(expr.FindAll(content, -1))
		if matches > maxMatches {
			maxMatches = matches
		}
		counts = append(counts, matches)
		return nil
	})
	if err != nil {
		logErr(err.Error())
	}
	return counts, maxMatches
}
