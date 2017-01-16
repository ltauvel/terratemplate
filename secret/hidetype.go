package secret

type HideType int
const (
	HIDE_NONE HideType = iota
	HIDE_PARTIAL
	HIDE_WHOLE
)