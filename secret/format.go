package secret

import (
	"regexp"
)

func formatSlice(text []string, hidetype HideType) []string {
	var result []string
	
	for _, t := range text {
		result = append(result, format(t, hidetype))
	}
	
	return result
}

func format(text string, hidetype HideType) string {
	
	if len(text) > 0 {
		var replacement = "********"

		r, _ := regexp.Compile("\\*\\(([^)]*)\\)")
		switch (hidetype){
			case HIDE_NONE: text = r.ReplaceAllString(text, "$1")
			case HIDE_PARTIAL: text = r.ReplaceAllString(text, replacement)
			case HIDE_WHOLE: text = replacement
		}
	}
	
	return text
}
