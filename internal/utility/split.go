package utility

import (
	"regexp"
	"strings"
)

func RegSplit(text string, regStr string) []string {
	reg := regexp.MustCompile(regStr)
	indexes := reg.FindAllStringIndex(text, -1)
	if len(indexes) == 0 {
		return []string{}
	}
	lastIndex := 0
	result := []string{}
	for _, element := range indexes {
		result = append(result, strings.TrimSpace(text[lastIndex:element[0]]))
		result = append(result, strings.TrimSpace(text[element[0]:element[1]]))
		lastIndex = element[1]
	}

	return append(result, strings.TrimSpace(text[lastIndex:]))
}
