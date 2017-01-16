package secret


func FormatSlice(text []string, hidetype HideType) []string {
	return formatSlice(text, hidetype)
}

func Format(text string, hidetype HideType) string {
	return format(text, hidetype)
}
