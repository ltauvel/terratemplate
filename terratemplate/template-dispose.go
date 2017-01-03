package main

import (
	"os"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
)

func DisposeTemplate(c *cli.Context) {

	if CheckPathExists(_instancepath) {

		destroyplan := _instancepath + "/destroy.plan"
		if CheckPathExists(destroyplan) {
		
			// Initilialize variables
			yellowText := color.New(color.FgYellow).Add(color.Underline)
			
			// Updating the terraform modules
			yellowText.Println("\r\nUpdating the terraform modules...\r\n* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * \r\n")
			ExecuteTerraform(_instancepath, "get", "-update", _templatepath)
			
			// Applying the terraform destruction plan
			yellowText.Println("\r\nApplying the terraform destruction plan...\r\n* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * \r\n")
			ExecuteTerraform(_instancepath, "apply", destroyplan)
			
			// Creating the backup directory
			os.MkdirAll(_backuppath,0777)
			os.RemoveAll(_backuppath)
			
			// Move the destroyed instance to the backups directory
			err := os.Rename(_instancepath, _backuppath)
			if err != nil {
				panic(err)
			}

		} else {
			panic("Cannot find the destruction plan at '" + destroyplan + "'")
		}
	} else {
		panic("Cannot find the specified instance at '" + _instancepath)
	}
}
