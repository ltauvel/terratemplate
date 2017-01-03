package main

import (
	"fmt"
	"os"
	"time"
	"strings"
	"regexp"
	"os/exec"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
    "path/filepath"
)


// Define global variables
/////////////////////////////////////////////////////////////////////////////////////////////////////////
var _executabledirectory string
var _workingdirectory string
var _debug bool
var _terraformpath string
var _templatepath string
var _instancepath string
var _backuppath string
var _extravars []string


// Startup function
/////////////////////////////////////////////////////////////////////////////////////////////////////////
func main() {
	app := cli.NewApp()
	app.Name = "terratemplate"
	app.Usage = "Provides template management functions with terraform development"
	app.Version = "1.0.0"
	app.Author = "Ludovic Tauvel"
	
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "debug mode",
		},
		cli.StringFlag{
			Name:  "templates-path, tp",
			Usage: "the `PATH` to the templates or [PWD/templates] or [%TERRATEMPLATE_BASEPATH%/templates] or ",
			EnvVar: "TERRATEMPLATE_TEMPLATEPATH",
		},
		cli.StringFlag{
			Name:  "instances-path, ip",
			Usage: "the `PATH` to the instances or [PWD/templates] or [%TERRATEMPLATE_BASEPATH%/instances] or ",
			EnvVar: "TERRATEMPLATE_INSTANCEPATH",
		},
		cli.StringFlag{
			Name:  "backup-path, bp",
			Usage: "the `PATH` to the backups or [PWD/templates] or [%TERRATEMPLATE_BASEPATH%/backups] or ",
			EnvVar: "TERRATEMPLATE_BACKUPPATH",
		},
	}
	
	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "template-name, tn",
			Usage: "the `NAME` of the template",
		},
		cli.StringFlag{
			Name:  "instance-id, ii",
			Usage: "the new template instance `ID`",
		},
	}
	
	app.Commands = []cli.Command{
		{
			Name:    "instanciate",
			Aliases: []string{"i"},
			Usage:   "instanciate the specified template",
			Action:  func(c *cli.Context) {
						InitializeGlobalVariables(c)
						InstanciateTemplate(c)
					 },
			Flags:	[]cli.Flag{
						cli.StringFlag{
							Name:  "template-name, tn",
							Usage: "the `NAME` of the template",
						},
						cli.StringFlag{
							Name:  "instance-id, ii",
							Usage: "the new template instance `ID`",
						},
						cli.StringSliceFlag{
							Name:  "var",
							Usage: "a `VAR` for the template. Can be used many times",
						},
						cli.BoolFlag{
							Name:  "force, f",
							Usage: "force the update even if the instance already exists",
						},
					},
		},
		{
			Name:    "dispose",
			Aliases: []string{"d"},
			Usage:   "dispose the specified template instance",
			Action:  func(c *cli.Context) {
						InitializeGlobalVariables(c)
						DisposeTemplate(c)
					 },
			Flags:	flags,
		},
	}
	
	app.Run(os.Args)
}



func CheckPathExists(path string) bool {
    _, err := os.Stat(path)
    if err == nil { return true }
    if os.IsNotExist(err) { return false }
    return true
}

func ExecuteTerraform(workingdir string, params ...string) {

	PrintDebug("command = " + _terraformpath, params...)

	// Launch terraform
	cmd := exec.Command(_terraformpath, params...)
	cmd.Dir = workingdir
	
	// Handle errot and standard output
	fmt.Println(string(GetValue(cmd.CombinedOutput()).([]byte)))
}

func GetValue(value interface{}, err interface{}) interface{} {
	if err != nil {
        panic(err)
    }
	return value
}

func PrintDebug(text string, args ...string) {
	debug := color.New(color.FgCyan).Add(color.Underline)
	if _debug {
		debug.Println("DEBUG: " + text + " " + strings.Join(args, " "))
	}
}

func InitializeGlobalVariables(c *cli.Context) {
	
	// Get the debug paameter value
	_debug = c.GlobalBool("debug")
	
	// Get the current executable path
	_executabledirectory = GetValue(filepath.Abs(filepath.Dir(os.Args[0]))).(string)
	PrintDebug("_executabledirectory =", _executabledirectory)
	
	// Get the working directory
	_workingdirectory = strings.Replace(GetValue(os.Getwd()).(string), " ", "\\ ", -1)
	PrintDebug("_workingdirectory =", _workingdirectory)
	
	// Get the terraform executable path
	_terraformpath = GetValue(exec.LookPath("terraform")).(string)
	PrintDebug("_terraformpath =", _terraformpath)
	
	// Compute the path to the template	
	_templatepath = c.GlobalString("templates-path")
	if len(_templatepath) == 0 {
		_templatepath = os.Getenv("TERRATEMPLATE_BASEPATH")
		if len(_templatepath) == 0 {
			_templatepath = _workingdirectory
		}
		_templatepath = _templatepath + "/templates"
	}
	if GetValue(regexp.MatchString("^\\.", _templatepath)) == true {
		_templatepath = filepath.Join(_workingdirectory, _templatepath)
	}
	_templatepath = filepath.Clean(_templatepath) + "/" +  c.String("template-name")
	PrintDebug("_templatepath =", _templatepath)

	// Compute the path to the instance
	_instancepath = c.GlobalString("instances-path")
	if len(_instancepath) == 0 {
		_instancepath = os.Getenv("TERRATEMPLATE_BASEPATH")
		if len(_instancepath) == 0 {
			_instancepath = _workingdirectory
		}
		_instancepath = _instancepath + "/instances"
	}
	if GetValue(regexp.MatchString("^\\.", _instancepath)) == true {
		_instancepath = filepath.Join(_workingdirectory, _instancepath)
	}
	_instancepath = filepath.Clean(_instancepath) + "/" +  c.String("template-name") + "/" + c.String("instance-id")
	PrintDebug("_instancepath =", _instancepath)
	
	// Compute the path to the backups
	_backuppath = c.GlobalString("backup-path")
	if len(_backuppath) == 0 {
		_backuppath = os.Getenv("TERRATEMPLATE_BASEPATH")
		if len(_backuppath) == 0 {
			_backuppath = _workingdirectory
		}
		_backuppath = _backuppath + "/backups"
	}
	if GetValue(regexp.MatchString("^\\.", _backuppath)) == true {
		_backuppath = filepath.Join(_workingdirectory, _backuppath)
	}
	_backuppath = filepath.Clean(_backuppath) + "/" +  c.String("template-name") + "/"+ c.String("instance-id") + "-" + time.Now().Local().Format("20060102150405")
	PrintDebug("_backuppath =", _backuppath)
	
	// Compute the extra vars
	for _, v := range c.StringSlice("var") {
		_extravars = append(_extravars, "-var")
		_extravars = append(_extravars, v)
	}
	PrintDebug("_extravars =", strings.Join(_extravars, " "))
}