package template

import (
	"time"
)

type Template struct {
	Handler ITemplateHandler
	Name string
	FullName string
	ModTime time.Time
}

func (template Template) HasChanged(comparepath string) bool {
	return template.Handler.HasChanged(template, comparepath)
}

func (template Template) Copy(destination string, force bool) Template { 
	return template.Handler.Copy(template, destination, force)
}
