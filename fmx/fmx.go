package fmx

import "unicode/utf8"

func Insert(subject string, content string) string {
	var index int

	for i, r := range subject {
		if r != ' ' && r != '\t' && r != '\n' {
			index = i
			break
		}
	}

	return subject[:index-(1+utf8.RuneCountInString(content))] + content + subject[index-1:]
}
