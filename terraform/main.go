package terraform

var TerraformPath string

func init() {
	TerraformPath = getTerraformPath()
}

func Exec(workingdir string, params ...string) {
	execute(workingdir, params...)
}