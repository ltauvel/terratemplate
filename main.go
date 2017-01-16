package main

import (
	"os"
	"github.com/codegangsta/cli"
	"github.com/ltauvel/gomodules/console"
)

// Startup function
/////////////////////////////////////////////////////////////////////////////////////////////////////////
func main() {
	app := cli.NewApp()
	app.Name = "terratemplate"
	app.Usage = "Provides template management functions with terraform development"
	app.Version = "1.0.0"
	app.Author = "Ludovic Tauvel"
	
	
	// Global arguments
	/////////////////////////////////////////////////////////////////////////////////////////////////////
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "debug mode",
		},
		cli.StringFlag{
			Name:  "templates-path, tp",
			Usage: "the `PATH` to the templates or [%PWD%/templates] or [%TERRATEMPLATE_BASEPATH%/templates] or ",
			EnvVar: "TERRATEMPLATE_TEMPLATEPATH",
		},
		cli.StringFlag{
			Name:  "instances-path, ip",
			Usage: "the `PATH` to the instances or [%PWD%/templates] or [%TERRATEMPLATE_BASEPATH%/instances] or ",
			EnvVar: "TERRATEMPLATE_INSTANCEPATH",
		},
		cli.StringFlag{
			Name:  "backups-path, bp",
			Usage: "the `PATH` to the backups or [%PWD%/templates] or [%TERRATEMPLATE_BASEPATH%/backups] or ",
			EnvVar: "TERRATEMPLATE_BACKUPPATH",
		},
	}
	
	app.Commands = []cli.Command{

		// Outputs commands
		/////////////////////////////////////////////////////////////////////////////////////////////////////
		{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "execute commands on the terraform output archives",
			Subcommands:	[]cli.Command{
				{
					Name:    "extract",
					Aliases: []string{"e"},
					Usage:   "extract a terratemplate output",
					Before:	 func(c *cli.Context) error {
								Intiallize(c)

								// Checking parameters
								if len(ttcontext.TemplateName) == 0 {
									console.PrintError("The template-name argument is mandatory.")
								}
								if len(ttcontext.InstanceId) == 0 {
									console.PrintError("The instance-id argument is mandatory.")
								}
								if len(ttcontext.DestinationPath) == 0 {
									console.PrintError("The path argument is mandatory.")
								}
								return nil
							 },
					Action:  func(c *cli.Context) {
								instance := TtInstanceManager{}
								instance.Extract()
							 },
					After:	 func(c *cli.Context) error {
								Release(c)
								return nil
							 },
					Flags:	[]cli.Flag{
								cli.StringFlag{
									Name:  "template-name, tn",
									Usage: "the `NAME` of the template",
								},
								cli.StringFlag{
									Name:  "instance-id, ii",
									Usage: "the `ID` of an existing instance",
								},
								cli.StringFlag{
									Name:  "path, p",
									Usage: "the output extraction destination `PATH`",
								},
								cli.StringFlag{
									Name:  "vault-password",
									Usage: "a `PASSWORD` that will be used to decrypt the terraform files.",
								},
								cli.StringFlag{
									Name:  "vault-password-file",
									Usage: "a `FILE` containing the password that will be used to decrypt the terraform files.",
								},
							},
				},
				{
					Name:    "encrypt",
					Aliases: []string{"e"},
					Usage:   "encrypt a terratemplate output",
					Before:	 func(c *cli.Context) error {
								Intiallize(c)

								// Checking parameters
								if len(ttcontext.TemplateName) == 0 {
									console.PrintError("The template-name argument is mandatory.")
								}
								if len(ttcontext.InstanceId) == 0 {
									console.PrintError("The instance-id argument is mandatory.")
								}
								if len(ttcontext.VaultPassword) == 0 {
									console.PrintError("The vault-password or vault-password-file argument is mandatory.")
								}
								
								return nil
							 },
					Action:  func(c *cli.Context) {
								instance := TtInstanceManager{}
								instance.Encrypt()
							 },
					After:	 func(c *cli.Context) error {
								Release(c)
								return nil
							 },
					Flags:	[]cli.Flag{
								cli.StringFlag{
									Name:  "template-name, tn",
									Usage: "the `NAME` of the template",
								},
								cli.StringFlag{
									Name:  "instance-id, ii",
									Usage: "the `ID` of an existing instance",
								},
								cli.StringFlag{
									Name:  "vault-password",
									Usage: "a `PASSWORD` that will be used to decrypt the terratemplate files",
								},
								cli.StringFlag{
									Name:  "vault-password-file",
									Usage: "a `FILE` containing the password that will be used to decrypt the terratemplate output",
								},
							},
				},
				{
					Name:    "decrypt",
					Aliases: []string{"e"},
					Usage:   "decrypt a terratemplate output",
					Before:	 func(c *cli.Context) error {
								Intiallize(c)

								// Checking parameters
								if len(ttcontext.TemplateName) == 0 {
									console.PrintError("The template-name argument is mandatory.")
								}
								if len(ttcontext.InstanceId) == 0 {
									console.PrintError("The instance-id argument is mandatory.")
								}
								if len(ttcontext.VaultPassword) == 0 {
									console.PrintError("The vault-password or vault-password-file argument is mandatory.")
								}
								
								return nil
							 },
					Action:  func(c *cli.Context) {
								instance := TtInstanceManager{}
								instance.Decrypt()
							 },
					After:	 func(c *cli.Context) error {
								Release(c)
								return nil
							 },
					Flags:	[]cli.Flag{
								cli.StringFlag{
									Name:  "template-name, tn",
									Usage: "the `NAME` of the template",
								},
								cli.StringFlag{
									Name:  "instance-id, ii",
									Usage: "the `ID` of an existing instance",
								},
								cli.StringFlag{
									Name:  "vault-password",
									Usage: "a `PASSWORD` that will be used to decrypt the terratemplate files",
								},
								cli.StringFlag{
									Name:  "vault-password-file",
									Usage: "a `FILE` containing the password that will be used to decrypt the terratemplate output",
								},
							},
				},
			},
		},
	
		// Templates commands
		/////////////////////////////////////////////////////////////////////////////////////////////////////
		{
			Name:    "template",
			Aliases: []string{"t"},
			Usage:   "execute commands on the templates",
			Subcommands:	[]cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list the templates",
					Before:	 func(c *cli.Context) error {
								Intiallize(c)
								return nil
							 },
					Action:  func(c *cli.Context) {
								template := TtTemplateManager{}
								template.List()
							 },
					After:	 func(c *cli.Context) error {
								Release(c)
								return nil
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
							},
				},
			},
		},
	
		// Instances arguments
		/////////////////////////////////////////////////////////////////////////////////////////////////////
		{
			Name:    "instance",
			Aliases: []string{"i"},
			Usage:   "execute commands on the instances",
			Subcommands:	[]cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list the instances",
					Before:	 func(c *cli.Context) error {
								Intiallize(c)
								
								// Checking parameters
								if ttcontext.ShowHistory && len(ttcontext.BackupsPath) == 0 {
									console.PrintError("Unable to determine the backups path.")
								}
								
								return nil
							 },
					Action:  func(c *cli.Context) {
								instance := TtInstanceManager{}
								instance.List()
							 },
					After:	 func(c *cli.Context) error {
								Release(c)
								return nil
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
								cli.BoolFlag{
									Name:  "show-history, sh",
									Usage: "show output history",
								},
							},
				},
				{
					Name:    "create",
					Aliases: []string{"c"},
					Usage:   "instanciate the specified template",
					Before:	 func(c *cli.Context) error {
								Intiallize(c)
								
								// Checking parameters
								if len(ttcontext.TemplateName) == 0 {
									console.PrintError("The template-name argument is mandatory.")
								}
								if len(ttcontext.InstanceId) == 0 {
									console.PrintError("The instance-id argument is mandatory.")
								}
								if len(ttcontext.TemplatesPath) == 0 {
									console.PrintError("Unable to determine the templates path.")
								}
								if len(ttcontext.InstancesPath) == 0 {
									console.PrintError("Unable to determine the instances path.")
								}
								if len(ttcontext.BackupsPath) == 0 {
									console.PrintError("Unable to determine the backups path.")
								}
								
								return nil
							 },
					Action:  func(c *cli.Context) {
								instance := TtInstanceManager{}
								instance.Create()
							 },
					After:	 func(c *cli.Context) error {
								Release(c)
								return nil
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
								cli.StringFlag{
									Name:  "vault-password",
									Usage: "a `PASSWORD` that will be used to decrypt the terraform files.",
								},
								cli.StringFlag{
									Name:  "vault-password-file",
									Usage: "a `FILE` containing the password that will be used to decrypt the terraform files.",
								},
							},
				},		
				{
					Name:    "dispose",
					Aliases: []string{"d"},
					Usage:   "dispose the specified template instance",
					Before:	 func(c *cli.Context) error {
								Intiallize(c)
								
								// Checking parameters
								if len(ttcontext.TemplateName) == 0 {
									console.PrintError("The template-name argument is mandatory.")
								}
								if len(ttcontext.InstanceId) == 0 {
									console.PrintError("The instance-id argument is mandatory.")
								}
								if len(ttcontext.TemplatesPath) == 0 {
									console.PrintError("Unable to determine the templates path.")
								}
								if len(ttcontext.InstancesPath) == 0 {
									console.PrintError("Unable to determine the instances path.")
								}
								if len(ttcontext.BackupsPath) == 0 {
									console.PrintError("Unable to determine the backups path.")
								}
								
								return nil
							 },
					Action:  func(c *cli.Context) {
								instance := TtInstanceManager{}
								instance.Dispose()
							 },
					After:	 func(c *cli.Context) error {
								Release(c)
								return nil
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
								cli.StringFlag{
									Name:  "vault-password",
									Usage: "a `PASSWORD` that will be used to decrypt the terraform files.",
								},
								cli.StringFlag{
									Name:  "vault-password-file",
									Usage: "a `FILE` containing the password that will be used to decrypt the terraform files.",
								},
							},
				},
			},
		},
	}
	
	app.Run(os.Args)
}


var ttcontext TtContext
func Intiallize(clicontext *cli.Context) {

	// Initializing console output writer
	console.Debug = clicontext.GlobalBool("debug")
	
	console.PrintDebug("Initializing ressources...")
	
	// Loading configuration file
	conf := TtConfiguration{}
	conf.Load()
	
	// Initializing context
	ttcontext = TtContext{}
	ttcontext.Initialize(clicontext, conf)

}

func Release(c *cli.Context) {

	console.PrintDebug("Releasing ressources...")

	// Removing temp directory
	os.RemoveAll(ttcontext.TempDir.FullName)
	
}
