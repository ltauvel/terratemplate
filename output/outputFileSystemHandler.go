package output

import (
	"strings"
	"github.com/ltauvel/gomodules/console"
	"github.com/ltauvel/gomodules/filesystem"
	"github.com/ltauvel/gomodules/targz"
)

type OutputFileSystemHandler struct {
	BasePath string
}

func (handler OutputFileSystemHandler) Type() string {
	return "FileSystem"
}

func (handler OutputFileSystemHandler) Extension() string {
	return ".tar.gz"
}

func (handler OutputFileSystemHandler) Get(templateName string, instanceId string) *Output {
	outputs := handler.List(templateName, instanceId)
	if len(outputs) > 0 {
		return outputs[0]
	} else {
		return nil
	}
}

func (handler OutputFileSystemHandler) List(templateName string, instanceId string) []*Output {
	result := []*Output{}
	
	// Looking for the templates in the directory defined as base path
	templates, _ := filesystem.ReadDir(handler.BasePath)
	for _, template := range templates {
		if len(templateName) == 0 || template.Name == templateName {
		
			// Looking for the instances in the current template directory	
			_, instances := template.Read()
			for _, instance := range instances {
				if len(instanceId) == 0 ||  instance.Name == instanceId + handler.Extension() {
					result = append(result, &Output{
						Handler: handler,
						TemplateName: template.Name,
						InstanceId: strings.Replace(instance.Name, handler.Extension(), "", -1),
						Name: instance.Name,
						FullName: instance.FullName,
						ModTime: instance.ModTime,
					})
				}
			}
		}
	}
	
	return result
}

func (handler OutputFileSystemHandler) Copy(source Output, destination string, force bool) Output {

	console.PrintDebug("Copying output")
	
	
	file := filesystem.LoadFile(source.FullName)
	file = file.Copy(destination, force)
	
	return Output{
			Handler: handler,
			TemplateName: source.TemplateName,
			InstanceId: strings.Replace(file.Name, handler.Extension(), "", -1),
			Name: file.Name,
			FullName: file.FullName,
			ModTime: file.ModTime,
		}

}

func (handler OutputFileSystemHandler) Create(sourcePath string, templateName string, outputPath string, passphrase string) Output {
	var result Output
	
	console.PrintDebug("Creating output")
	
	// Compressing all the output files
	targz.CompressFolder(outputPath, sourcePath)
	
	// Loading the created output
	file := filesystem.LoadFile(outputPath)
	result = Output{
			Handler: handler,
			TemplateName: templateName,
			InstanceId: strings.Replace(file.Name, handler.Extension(), "", -1),
			Name: file.Name,
			FullName: file.FullName,
			ModTime: file.ModTime,
		}
	
	// Encrypt the archive if needed
	if len(passphrase) > 0 {
		if !result.Encrypt(passphrase) {
			console.PrintError("Cannot encrypt to the specified output archive.")
		}
	}
	
	return result
}