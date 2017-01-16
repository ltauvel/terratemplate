package template

import (
	"github.com/ltauvel/gomodules/console"
	"github.com/ltauvel/gomodules/filesystem"
)

type TemplateFileSystemHandler struct {
	BasePath string
}

func (handler TemplateFileSystemHandler) Type() string {
	return "FileSystem"
}

func (handler TemplateFileSystemHandler) Get(templateName string) *Template {
	templates := handler.List(templateName)
	if len(templates) > 0 {
		return templates[0]
	} else {
		return nil
	}
}

func (handler TemplateFileSystemHandler) List(templateName string) []*Template {
	result := []*Template{}
	
	// Looking for the templates in the directory defined as base path
	templates, _ := filesystem.ReadDir(handler.BasePath)
	for _, template := range templates {
		if len(templateName) == 0 || template.Name == templateName {
		
			result = append(result, &Template{
				Handler: handler,
				Name: template.Name,
				FullName: template.FullName,
				ModTime: template.ModTime,
			})
		}
	}
	
	return result
}

func (handler TemplateFileSystemHandler) Copy(source Template, destination string, force bool) Template {
	
	console.PrintDebug("Copying template")
	
	dir := filesystem.LoadDir(source.FullName)
	dir = dir.Copy(destination, true, force)
	
	return Template{
			Handler: handler,
			Name: dir.Name,
			FullName: dir.FullName,
			ModTime: dir.ModTime,
		}
}
