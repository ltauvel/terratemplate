package main

import (
	"os"
	"time"
	"strings"
	"strconv"
	"io/ioutil"
	"github.com/codegangsta/cli"
	"github.com/ltauvel/gomodules/console"
	"github.com/ltauvel/gomodules/filesystem"
	"github.com/ltauvel/terratemplate/output"
	"github.com/ltauvel/terratemplate/template"
	"github.com/ltauvel/terratemplate/secret"
)

type TtContext struct {
	Debug bool
	TempDir filesystem.Directory
	WorkingDir string
	Force bool
	ShowHistory bool
	DestinationPath string
	TemplatesPath string
	TemplatePath string
	InstancesPath string
	InstancePath string
	BackupsPath string
	BackupPath string
	OutputName string
	OutputPath string
	OutputBackupPath string
	TemplateName string
	InstanceId string
	ExtraVars []string
	VaultPassword string
	VaultPasswordFile string
	OutputHandler output.IOutputHandler
	TemplateHandler template.ITemplateHandler
}

func (ttcontext *TtContext) Initialize(clicontext *cli.Context, config TtConfiguration) {

	// Get the debug parameter value
	ttcontext.Debug = clicontext.GlobalBool("debug")
	
	console.PrintDebugSection("Context")
	
	console.PrintDebug("Debug =", strconv.FormatBool(ttcontext.Debug))

	// Get the show-history parameter value
	ttcontext.ShowHistory = clicontext.Bool("show-history")
	console.PrintDebug("ShowHistory = ", strconv.FormatBool(ttcontext.ShowHistory))
	
	// Set a temporary directory
	tempdir, _ := ioutil.TempDir("", "tf")
	ttcontext.TempDir = filesystem.LoadDir(tempdir)
	console.PrintDebug("TempDir = ", ttcontext.TempDir.FullName)
	
	// Get the working directory
	ttcontext.WorkingDir = GetValue(os.Getwd()).(string)
	console.PrintDebug("WorkingDir = ", ttcontext.WorkingDir)
	
	// Get the force parameter value
	ttcontext.Force = clicontext.Bool("force")
	console.PrintDebug("Force =", strconv.FormatBool(ttcontext.Force))
	
	// Get the template-name parameter value
	ttcontext.TemplateName = clicontext.String("template-name")
	console.PrintDebug("TemplateName =", ttcontext.TemplateName)

	// Get the instance-id parameter value
	ttcontext.InstanceId = clicontext.String("instance-id")
	console.PrintDebug("InstanceId =", ttcontext.InstanceId)
	
	// Compute the path to the templates	
	ttcontext.TemplatesPath = filesystem.GetValidPath(clicontext.GlobalString("templates-path"), config.TemplatesPath, "%TERRATEMPLATE_BASEPATH%/templates", "./templates")
	console.PrintDebug("TemplatesPath =", ttcontext.TemplatesPath)
	
	// Compute the path to the template
	if len(ttcontext.TemplatesPath) > 0 && len(ttcontext.TemplateName) > 0 {
		ttcontext.TemplatePath = filesystem.JoinPath(ttcontext.TemplatesPath, ttcontext.TemplateName)
	}
	console.PrintDebug("TemplatePath =", ttcontext.TemplatePath)

	// Compute the path to the instances
	ttcontext.InstancesPath = filesystem.GetValidPath(clicontext.GlobalString("instances-path"), config.InstancesPath, "%TERRATEMPLATE_BASEPATH%/instances", "./instances")
	console.PrintDebug("InstancesPath =", ttcontext.InstancesPath)
	
	// Compute the path to the instance
	if len(ttcontext.InstancesPath) > 0 && len(ttcontext.TemplateName) > 0 {
		ttcontext.InstancePath = filesystem.JoinPath(ttcontext.InstancesPath, ttcontext.TemplateName)
	}
	console.PrintDebug("InstancePath =", ttcontext.InstancePath)
	
	// Compute the path to the backups	
	ttcontext.BackupsPath = filesystem.GetValidPath(clicontext.GlobalString("backups-path"), config.BackupsPath, "%TERRATEMPLATE_BASEPATH%/backups", "./backups")
	console.PrintDebug("BackupsPath =", ttcontext.BackupsPath)
	
	// Compute the path to the backup
	if len(ttcontext.BackupsPath) > 0 && len(ttcontext.TemplateName) > 0 && len(ttcontext.InstanceId) > 0 {
		ttcontext.BackupPath = filesystem.JoinPath(ttcontext.BackupsPath, ttcontext.TemplateName, ttcontext.InstanceId)
	}
	console.PrintDebug("BackupPath =", ttcontext.BackupPath)
	
	// Compute the output name
	if len(ttcontext.InstanceId) > 0 {
		ttcontext.OutputName = ttcontext.InstanceId + ".tar.gz"
	}
	console.PrintDebug("OutputName = ", ttcontext.OutputName)
	
	// Set the current output archive path
	if len(ttcontext.InstancePath) > 0 && len(ttcontext.InstanceId) > 0 && len(ttcontext.OutputName) > 0 {
		ttcontext.OutputPath = filesystem.JoinPath(ttcontext.InstancePath, ttcontext.OutputName)
	}
	console.PrintDebug("OutputPath = ", ttcontext.OutputPath)
	
	// Compute the path to the output backup
	if len(ttcontext.BackupPath) > 0  {
		ttcontext.OutputBackupPath = filesystem.JoinPath(ttcontext.BackupPath, time.Now().Local().Format("20060102150405") + ".tar.gz")
	}
	console.PrintDebug("OutputBackupPath =", ttcontext.OutputBackupPath)
	
	// Compute the destination path
	ttcontext.DestinationPath = clicontext.String("path")
	console.PrintDebug("DestinationPath =", ttcontext.DestinationPath)
	
	// Compute the extra vars
	for _, v := range clicontext.StringSlice("var") {
		ttcontext.ExtraVars = append(ttcontext.ExtraVars, "-var")
		ttcontext.ExtraVars = append(ttcontext.ExtraVars, v)
	}
	console.PrintDebug("ExtraVars =", secret.Format(strings.Join(ttcontext.ExtraVars, " "), secret.HIDE_PARTIAL))
	
	// Get the vault password
	ttcontext.VaultPassword = clicontext.String("vault-password")
	if len(ttcontext.VaultPassword) == 0 {
		ttcontext.VaultPasswordFile = filesystem.GetValidPath(clicontext.String("vault-password-file"))
		if len(ttcontext.VaultPasswordFile) > 0 {
			console.PrintDebug("VaultPasswordFile =", ttcontext.VaultPasswordFile)
			content, err := ioutil.ReadFile(ttcontext.VaultPasswordFile)
			if err != nil {
				panic(err.Error())
			}
			ttcontext.VaultPassword = string(content)
		}
	}
	console.PrintDebug("VaultPassword =", secret.Format(ttcontext.VaultPassword, secret.HIDE_WHOLE))

	
	// Instanciate the output Handler
	ttcontext.OutputHandler = output.OutputFileSystemHandler{ BasePath: ttcontext.InstancesPath }
	console.PrintDebug("OutputHandler =", ttcontext.OutputHandler.Type())
	
	// Instanciate the template Handler
	ttcontext.TemplateHandler = template.TemplateFileSystemHandler{ BasePath: ttcontext.TemplatesPath }
	console.PrintDebug("TemplateHandler =", ttcontext.TemplateHandler.Type())
	
	console.PrintDebugSection("")
}
