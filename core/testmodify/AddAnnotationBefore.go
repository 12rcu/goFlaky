package testmodify

import "regexp"

func AddAnnotationBefore(content string, regex *regexp.Regexp, insertText func(matchNum int) string) string {
	counter := 0
	return regex.ReplaceAllStringFunc(content, func(s string) string {
		result := insertText(counter) + "\n" + s
		counter++
		return result
	})
}
