package main

import (
	"os"
	"time"
	"github.com/ltauvel/gomodules/filesystem"
	"github.com/ltauvel/gomodules/console"
	"github.com/ltauvel/terratemplate/terraform"
)

type TtInstanceManager struct {
}


func (instance TtInstanceManager) List() {

	var content [][]string

	// Generate the content table
	for _, o := range ttcontext.OutputHandler.List(ttcontext.TemplateName, ttcontext.InstanceId) {
		content = append(content, []string{o.TemplateName, o.InstanceId, o.ModTime.Format("2006-01-02 15:04:05"), o.FullName})
	}
	
	// Print the table
	console.PrintTable([]string{"TEMPLATE NAME", "INSTANCE ID", "MODIFICATION DATE", "OUTPUT PATH"}, content, 2)
}

func (instance TtInstanceManager) Extract() {

	// Load the previous instance output
	ttoutput := ttcontext.OutputHandler.Get(ttcontext.TemplateName, ttcontext.InstanceId)
	if ttoutput != nil {
	
		ttoutput.Extract(ttcontext.DestinationPath, ttcontext.VaultPassword)

	} else {
		console.PrintError("Cannot find the output for the specified instance " + ttcontext.OutputPath)
	}
}

func (instance TtInstanceManager) Encrypt() {

	// Load the previous instance output
	ttoutput := ttcontext.OutputHandler.Get(ttcontext.TemplateName, ttcontext.InstanceId)
	if ttoutput != nil {
	
		if !ttoutput.IsEncrypted() {
			ttoutput.Encrypt(ttcontext.VaultPassword)
		} else {
			console.PrintError("The specified instance output has already been encrypted. Please decrypt the output first.")
		}

	} else {
		console.PrintError("Cannot find the output for the specified instance " + ttcontext.OutputPath)
	}
}

func (instance TtInstanceManager) Decrypt() {

	// Load the previous instance output
	ttoutput := ttcontext.OutputHandler.Get(ttcontext.TemplateName, ttcontext.InstanceId)
	if ttoutput != nil {
	
		if ttoutput.IsEncrypted() {
			ttoutput.Decrypt(ttcontext.VaultPassword)
		} else {
			console.PrintError("The specified instance output has not been encrypted.")
		}

	} else {
		console.PrintError("Cannot find the output for the specified instance " + ttcontext.OutputPath)
	}
}

func (instance TtInstanceManager) Create() {

	// Load the template and check that it exists
	template := ttcontext.TemplateHandler.Get(ttcontext.TemplateName)
	if template != nil {
	
		// Load the previous instance output
		ttoutput := ttcontext.OutputHandler.Get(ttcontext.TemplateName, ttcontext.InstanceId)
		if ttcontext.Force || ttoutput == nil {
		
			// Creating the execution directory
			execdir := ttcontext.TempDir.Create("latest")

			// Copying template source to the execution directory
			template.Copy(execdir.FullName, true)
		
			// If an output already exists then extract it to the temp directory
			if ttoutput != nil {
				console.PrintSection("Reading the terraform outputs archive for the current instance...")
				ttoutput.Extract(ttcontext.TempDir.FullName, ttcontext.VaultPassword)
			}
			
			// Updating the terraform modules
			console.PrintSection("Updating the terraform modules...")
			terraform.Exec(execdir.FullName, "get", "-update", execdir.FullName)
			
			// Generting the terraform execution plan
			console.PrintSection("Generating the terraform execution plan...")
			terraform.Exec(execdir.FullName, JoinStringSlices([]string{"plan", "-var", "instance_id=" + ttcontext.InstanceId}, ttcontext.ExtraVars, []string{"-out", execdir.FullName + "/apply.plan", execdir.FullName})...)

			// Applying the terraform execution plan
			console.PrintSection("Applying the terraform execution plan...")
			terraform.Exec(execdir.FullName, "apply", execdir.FullName + "/apply.plan")

			// Generting the terraform destruction plan
			console.PrintSection("Generating the terraform destruction plan...")
			terraform.Exec(execdir.FullName, JoinStringSlices([]string{"plan", "-destroy", "-var", "instance_id=" + ttcontext.InstanceId}, ttcontext.ExtraVars, []string{"-out", execdir.FullName + "/destroy.plan", execdir.FullName})...)

			// Generting the terraform outputs archive
			console.PrintSection("Generating the terraform outputs archive...")
			execdir.Copy(filesystem.JoinPath(ttcontext.TempDir.FullName, time.Now().Local().Format("20060102150405")), true, true)
			os.MkdirAll(ttcontext.InstancePath, 0777)
			ttcontext.OutputHandler.Create(ttcontext.TempDir.FullName, ttcontext.TemplateName, ttcontext.OutputPath, ttcontext.VaultPassword)
				
		} else {
			console.PrintError("An instance with the same id already exists at '" + ttoutput.FullName + "'. Use the update command instead if you want to update the existing instance.")
		}
	} else {
		console.PrintError("The specified template cannot be found at '" + ttcontext.TemplatePath + "'")
	}
}

func (instance TtInstanceManager) Update() {

	// Load the template and check that it exists
	template := ttcontext.TemplateHandler.Get(ttcontext.TemplateName)
	if template != nil {
	
		// Load the previous instance output
		ttoutput := ttcontext.OutputHandler.Get(ttcontext.TemplateName, ttcontext.InstanceId)
		if ttoutput != nil {
						
			// Creating the execution directory
			execdir := ttcontext.TempDir.Create("latest")

			// Copying template source to the execution directory
			template.Copy(execdir.FullName, true)
		
			// Extract the previous output to the temp directory
			console.PrintSection("Reading the terraform outputs archive for the current instance...")
			ttoutput.Extract(ttcontext.TempDir.FullName, ttcontext.VaultPassword)
			
			// Check if template has changes since last run
			if ttcontext.Force || template.HasChanged(execdir.FullName) {

				// Updating the terraform modules
				console.PrintSection("Updating the terraform modules...")
				terraform.Exec(execdir.FullName, "get", "-update", execdir.FullName)
				
				// Generting the terraform execution plan
				console.PrintSection("Generating the terraform execution plan...")
				terraform.Exec(execdir.FullName, JoinStringSlices([]string{"plan", "-var", "instance_id=" + ttcontext.InstanceId}, ttcontext.ExtraVars, []string{"-out", execdir.FullName + "/apply.plan", execdir.FullName})...)

				// Applying the terraform execution plan
				console.PrintSection("Applying the terraform execution plan...")
				terraform.Exec(execdir.FullName, "apply", execdir.FullName + "/apply.plan")

				// Generting the terraform destruction plan
				console.PrintSection("Generating the terraform destruction plan...")
				terraform.Exec(execdir.FullName, JoinStringSlices([]string{"plan", "-destroy", "-var", "instance_id=" + ttcontext.InstanceId}, ttcontext.ExtraVars, []string{"-out", execdir.FullName + "/destroy.plan", execdir.FullName})...)

				// Generting the terraform outputs archive
				console.PrintSection("Generating the terraform outputs archive...")
				execdir.Copy(filesystem.JoinPath(ttcontext.TempDir.FullName, time.Now().Local().Format("20060102150405")), true, true)
				os.MkdirAll(ttcontext.InstancePath, 0777)
				ttcontext.OutputHandler.Create(ttcontext.TempDir.FullName, ttcontext.TemplateName, ttcontext.OutputPath, ttcontext.VaultPassword)
				
			} else {
				console.PrintError("No changes have been detected in the template since the last execution. Use the 'force' argument if you wish to execute it anyway.")
			}
		} else {
			console.PrintError("Cannot find any existing instance at '" + ttoutput.FullName + "'. Use the 'create' command if you wish to create a new instance.")
		}
	} else {
		console.PrintError("The specified template cannot be found at '" + ttcontext.TemplatePath + "'")
	}
}

func (instance TtInstanceManager) Dispose() {

	// Load the previous instance output and ceck that it exists
	ttoutput := ttcontext.OutputHandler.Get(ttcontext.TemplateName, ttcontext.InstanceId)
	if ttoutput != nil {
	
		// Creating the execution directory
		execdir := ttcontext.TempDir.Create("latest")
	
		// If an output already exists then extract it to the temp directory
		console.PrintSection("Reading the terraform outputs archive for the current instance...")
		ttoutput.Extract(ttcontext.TempDir.FullName, ttcontext.VaultPassword)

		destroyplan := filesystem.JoinPath(execdir.FullName, "destroy.plan")
		if filesystem.CheckPathExists(destroyplan) {
		
			// Updating the terraform modules
			console.PrintSection("Updating the terraform modules...")
			terraform.Exec(ttcontext.TempDir.FullName, "get", "-update", ttcontext.TemplatePath)
			
			// Applying the terraform destruction plan
			console.PrintSection("Applying the terraform destruction plan...")
			terraform.Exec(ttcontext.TempDir.FullName, "apply", destroyplan)
			
			// Generting the terraform outputs archive
			console.PrintSection("Generating the terraform outputs archive...")
			execdir.Copy(filesystem.JoinPath(ttcontext.TempDir.FullName, time.Now().Local().Format("20060102150405")), true, true)
			os.MkdirAll(ttcontext.InstancePath, 0777)	
			ttcontext.OutputHandler.Create(ttcontext.TempDir.FullName, ttcontext.TemplateName, ttcontext.OutputPath, ttcontext.VaultPassword)
			
			console.PrintSection("Backuping the terraform output...")
			
			// Creating the backup directory
			os.MkdirAll(ttcontext.BackupPath,0777)
			
			// Move the destroyed instance to the backups directory
			err := os.Rename(ttcontext.OutputPath, ttcontext.OutputBackupPath)
			if err != nil {
				panic(err)
			}

		} else {
			console.PrintError("Cannot find the destruction plan" + destroyplan)
		}
	} else {
		console.PrintError("Cannot find the specified instance " + ttcontext.OutputPath)
	}
}
