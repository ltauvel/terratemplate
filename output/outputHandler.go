package output

type IOutputHandler interface {
	Type() string
	List(templateName string, instanceId string) []*Output
	Get(templateName string, instanceId string) *Output
	Create(sourcePath string, templateName string, outputPath string, passphrase string) Output
	Copy(source Output, destination string, force bool) Output
}
