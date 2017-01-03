package main

import (
	"os"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
)

func JoinSlices(params ...[]string) []string {
	var result = []string{}
	for _, p := range params {
		for _, v := range p {
			result = append(result, v)
		}
	}
	return result
}

func InstanciateTemplate(c *cli.Context) {

	if CheckPathExists(_templatepath) {
	
		if c.Bool("force") || !CheckPathExists(_instancepath) {
		
			// Initilialize variables
			yellowText := color.New(color.FgYellow).Add(color.Underline)

			// Creating the instance directory
			os.MkdirAll(_instancepath, 0777)

			// Updating the terraform modules
			yellowText.Println("\r\nUpdating the terraform modules...\r\n* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * \r\n")
			ExecuteTerraform(_instancepath, "get", "-update", _templatepath)
			
			// Generting the terraform execution plan
			yellowText.Println("\r\nGenerating the terraform execution plan...\r\n* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * \r\n")
			ExecuteTerraform(_instancepath, JoinSlices([]string{"plan", "-var", "instance_id=" + c.String("instance-id")}, _extravars, []string{"-out", _instancepath + "/apply.plan", _templatepath})...)

			// Applying the terraform execution plan
			yellowText.Println("\r\nApplying the terraform execution plan...\r\n* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * \r\n")
			ExecuteTerraform(_instancepath, "apply", _instancepath + "/apply.plan")

			// Generting the terraform destruction plan
			yellowText.Println("\r\nGenerating the terraform destruction plan...\r\n* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * \r\n")
			ExecuteTerraform(_instancepath, JoinSlices([]string{"plan", "-destroy", "-var", "instance_id=" + c.String("instance-id")}, _extravars, []string{"-out", _instancepath + "/destroy.plan", _templatepath})...)

		} else {
			panic("An instance with the same id already exists at '" + _instancepath + "'. Use the update command instead if you want to update the existing instance.")
		}
	} else {
		panic("The specified template cannot be found at '" + _templatepath + "'")
	}
}