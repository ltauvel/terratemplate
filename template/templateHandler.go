package template

type ITemplateHandler interface {
	Type() string
	Get(templateName string) *Template
	List(templateName string) []*Template
	Copy(source Template, destination string, force bool) Template
	HasChanged(source Template, comparepath string) bool
}