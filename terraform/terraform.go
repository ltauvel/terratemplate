package terraform

import (
	"os"	
	"os/exec"
	"github.com/ltauvel/gomodules/console"
	"github.com/ltauvel/terratemplate/secret"
)

// Run the terraform executable
func execute(workingdir string, params ...string) {

	console.PrintDebug("command = " + TerraformPath, secret.FormatSlice(params, secret.HIDE_PARTIAL)...)

	// Launch terraform
	cmd := exec.Command(TerraformPath, secret.FormatSlice(params, secret.HIDE_NONE)...)
	cmd.Dir = workingdir
	cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
    cmd.Run()
	
}

