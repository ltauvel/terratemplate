package main

import (
	"github.com/ltauvel/gomodules/console"
)

type TtTemplateManager struct {
}

func (instance TtTemplateManager) List() {

	var content [][]string

	// Generate the content table
	for _, t := range ttcontext.TemplateHandler.List(ttcontext.TemplateName) {
		content = append(content, []string{t.Name, t.ModTime.Format("2006-01-02 15:04:05"), t.FullName})
	}
	
	// Print the table
	console.PrintTable([]string{"TEMPLATE NAME", "MODIFICATION DATE", "TEMPLATE PATH"}, content, 2)
	
}
