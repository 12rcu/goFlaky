package testmodify

import "regexp"

func AddFirstAnnotationBefore(content string, regex *regexp.Regexp, insertText string) string {
	replaced := false
	return regex.ReplaceAllStringFunc(content, func(s string) string {
		if replaced {
			return s
		}
		result := insertText + "\n" + s
		replaced = true
		return result
	})
}
