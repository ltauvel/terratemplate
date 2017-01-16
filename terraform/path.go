package terraform


import (
	"os"	
	"os/exec"
    "path/filepath"
	"github.com/ltauvel/gomodules/console"
	"github.com/ltauvel/gomodules/filesystem"
)


// Return the terraform executable path
func getTerraformPath() string {

	// Look if the tarraform exeutable is in the path environment variable
	terraformpath, err := exec.LookPath("terraform")
	
	// If the terraform executable has not been foud in %PATH%
	// Try to locate it in terratemplate executable directory or in working directory
	if err != nil {
		tfpath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		terraformpath = filesystem.GetValidPath( filesystem.JoinPath(tfpath, "terraform.exe"), "./terraform.exe" )
	}
	
	if len(terraformpath) == 0 {
		console.PrintError("Unable to find terraform executable.")
	}

	return terraformpath
	
}