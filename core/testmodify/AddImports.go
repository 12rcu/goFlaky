package testmodify

import "regexp"

func AddImports(content string, insertText string, regex *regexp.Regexp) string {
	if len(regex.FindString(content)) == 0 {
		return insertText + "\n" + content
	}
	replaced := false
	return regex.ReplaceAllStringFunc(content, func(s string) string {
		if replaced {
			return s
		}
		result := s + "\n" + insertText + "\n"
		replaced = true
		return result
	})
}
