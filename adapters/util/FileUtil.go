package util

import (
	"log"
	"os"
	"path/filepath"
)

func SearchFile(path string, expr func(path string, fileName string) bool) []string {
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
		log.Default().Print(err)
	}
	return files
}

func SearchFileWithContent(path string, expr func(path string, fileName string, fileContent string) bool) []string {
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
		log.Default().Print(err)
	}
	return files
}
