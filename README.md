# Terratemplate
===

`terratemplate` is as a command line utility written in [Go](https://github.com/golang/go) and is aimed at providing template management and output encryption when working with [Terraform](https://www.terraform.io). 



<!-- toc -->
- [Overview](#overview)
- [Object Types](#objecttypes)
  * [Template](#template)
  * [Instance](#instance)
  * [Output](#output)
  * [Backup](#backup)
- [Installation](#install)
  * [Supported platforms](#supported-platforms)
  * [Install notes](#install-notes)
  * [Configuration](#config)
  * [Configuration File](#config-file)
- [Getting started](#get-start)
  * [Encryption](#encryption)
  * [Template Commands](#template-commands)
  * [Instance Commands](#instance-commands)
  * [Output Commands](#output-commands)

<!-- tocstop -->

## Overview

Certain type of environment need to be duplicated many times like DEV or QA environments. `Terraform` can handle this either by duplicating the tf files in many folders and changing the name of each resource to create or by using tf files waiting for input variables which will be used to uniquely name the resources. In both previous solution a separate folder should be created to store the tfsstate files.

`Terratemplate` is aimed at handling this by using variabilized tf files. It additionally comes with an output encryption functionality that let you vault all the `Terraform` output using a passphrase. 

## Object Types

### Template

A template is a folder containing tf files used to create resources. All the templates need to be in a same parent folder.

	..
	|_ templates
		|_ template_1
		|	|_ tf_files
		|	|_ ...
		|_ template2
			|_ tf_files
			|_ ...
			
### Instance

An instance is the result of a template execution. All the instances are physically stored like this:

	..
	|_ instances
		|_ template_1
		|	|_ instance_1.tar.gz
		|	|_ instance_2.tar.gz
		|	|_ ...
		|_ template2
			|_ instance_1.tar.gz
			|_ instance_2.tar.gz
			|_ ...

### Output

An output is a tar.gz archive which contains a template instance outputs. See [Instance](#instance). It contains all the `Terraform` runs history.

	instance_1.tar.gz
	|_ date_1
	|	|_ tf_modules
	|	|_ tf_files
	|	|_ tfstate_files
	|	|_ execution plan
	|	|_ destruction plan
	|_ date_2
	|	|_ tf_modules
	|	|_ tf_files
	|	|_ tfstate_files
	|	|_ execution plan
	|	|_ destruction plan
	|_ ...
	|_ latest
		|_ tf_modules
		|_ tf_files
		|_ tfstate_files
		|_ execution plan
		|_ destruction plan
			
### Backup

All the disposed instance are automatically placed in a backup folder to preserve instances history. They are stored like this:

	..
	|_ backups
		|_ template_1
		|	|_ instance_1.dipose_date.tar.gz
		|	|_ instance_2.dipose_date.tar.gz
		|	|_ ...
		|_ template2
			|_ instance_1.dipose_date.tar.gz
			|_ instance_2.dipose_date.tar.gz
			|_ ...
			
			
## Installation

### Supported platforms

`Terratemplate` have been tested and compiled for Windows at this time. But more operating system will be supported soon.

### Install notes

First of all, `Terraform` needs to be installed the way you prefer. Then you have to download `Terratemplate` binary.
`Terratemplate` will look for `Terraform` binary in the PATH, in the same folder than the `Terratemplate` binary or in the current working directory.

### Configuration

`Terratemplate` has to find all the path to the special folders listed above. You can specify these special folders this way:

Path to the templates:

	- Using the following command line argument: terratemplate --templates-path "PATH"
	- In a folder named "templates" in the working directory: %PWD%/templates
	- In a folder named "templates" in the directory specified in a TERRATEMPLATE_BASEPATH environment variable: %TERRATEMPLATE_BASEPATH%/templates
	- In a directory specified in a %TERRATEMPLATE_TEMPLATEPATH% environment variable

Path to the instances:

	- Using the following command line argument: terratemplate --instances-path "PATH"
	- In a folder named "instances" in the working directory: %PWD%/instances
	- In a folder named "instances" in the directory specified in a TERRATEMPLATE_BASEPATH environment variable: %TERRATEMPLATE_BASEPATH%/instances
	- In a directory specified in a %TERRATEMPLATE_TEMPLATEPATH% environment variable

Path to the backups:

	- Using the following command line argument: terratemplate --backups-path "PATH"
	- In a folder named "backups" in the working directory: %PWD%/backups
	- In a folder named "backups" in the directory specified in a TERRATEMPLATE_BASEPATH environment variable: %TERRATEMPLATE_BASEPATH%/backups
	- In a directory specified in a %TERRATEMPLATE_TEMPLATEPATH% environment variable

### Configuration File

You can alternatively use a `terratemplate.cfg` JSON configuration file to specify the three paths above. `Terratemplate` will look for it in the working directory and pick the paths.

	{
		"templates-path": "./templates",
		"instances-path": "./instances",
		"backups-path": "./backups",
	}

## Getting started

### Encryption

`terratemplate` support encryption/decryption in all the commands. To encrypt or decrypt either use the `--vault-password "*****"` argument or the `--vault-password-file ".\xxx"` which must specify the path to a file containing the encryption passphrase.

### Template Commands

To get help type the command `terratemplate template -h`
	
	NAME:
	terratemplate template - execute commands on the templates

	USAGE:
	   terratemplate template command [command options] [arguments...]

	COMMANDS:
		 list, l  list the templates

	OPTIONS:
	   --help, -h  show help

* Listing all the templates with command `terratemplate template list`
	
		TEMPLATE NAME  MODIFICATION DATE    TEMPLATE PATH
		-------------  -------------------  ------------------
		DEV            2016-12-20 15:14:15  \tmp\templates\DEV
		QA             2016-12-20 15:14:15  \tmp\templates\QA
		-------------  -------------------  ------------------
		TOTAL: 2

### Instances Commands

To get help type the command `terratemplate template -h`

	NAME:
	   terratemplate instance - execute commands on the instances

	USAGE:
	   terratemplate instance command [command options] [arguments...]

	COMMANDS:
		 list, l     list the instances
		 create, c   instantiate the specified template
		 dispose, d  dispose the specified template instance

	OPTIONS:
	   --help, -h  show help

* Listing all the instances with command `terratemplate instance list`

		TEMPLATE NAME  INSTANCE ID  MODIFICATION DATE    OUTPUT PATH
		-------------  -----------  -------------------  ---------------------------
		DEV            1            2017-01-17 13:39:04  \tmp\instances\DEV\1.tar.gz
		DEV            2            2017-01-17 13:39:12  \tmp\instances\DEV\2.tar.gz
		QA             1            2017-01-17 13:38:28  \tmp\instances\QA\1.tar.gz
		-------------  -----------  -------------------  ---------------------------
		TOTAL: 3
	
* Creating an instance from a template with command `terratemplate instance create --template-name QA --instance-id 1 --var "subscription_id=***" -
-var "client_id=***" --var "client_secret=***" --var "tenant_id=***"`

		Updating the terraform modules...
		* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
		...

		Generating the terraform execution plan...
		* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *

		Path: \tmp\tf815291415\latest/apply.plan

		+ azurerm_resource_group.rg
			location:         "westeurope"
			name:             "r71rg"
			tags.%:           "1"
			tags.environment: "qa1"

		Plan: 1 to add, 0 to change, 0 to destroy.

		Applying the terraform execution plan...
		* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *

		azurerm_resource_group.rg: Creating...
		  location:         "" => "westeurope"
		  name:             "" => "r71rg"
		  tags.%:           "" => "1"
		  tags.environment: "" => "qa1"
		azurerm_resource_group.rg: Creation complete

		Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

		State path: terraform.tfstate

		Generating the terraform destruction plan...
		* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *

		Path: \tmp\tf815291415\latest/destroy.plan

		- azurerm_resource_group.rg


		Plan: 0 to add, 0 to change, 1 to destroy.

		Generating the terraform outputs archive...
		* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
		...

* Updating an existing instance using the `--force` argument: `terratemplate instance create --template-name QA --instance-id 1 --var "subscription_id=***" -
-var "client_id=***" --var "client_secret=***" --var "tenant_id=***" --force`
	
* Removing an instance with command `terratemplate instance dispose --template-name QA --instance-id 1`

		Reading the terraform outputs archive for the current instance...
		* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
		...
		
		Updating the terraform modules...
		* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
		...
		
		Applying the terraform destruction plan...
		* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *

		azurerm_resource_group.rg: Destroying...
		azurerm_resource_group.rg: Still destroying... (10s elapsed)
		azurerm_resource_group.rg: Still destroying... (20s elapsed)
		azurerm_resource_group.rg: Still destroying... (30s elapsed)
		azurerm_resource_group.rg: Still destroying... (40s elapsed)
		azurerm_resource_group.rg: Destruction complete

		Apply complete! Resources: 0 added, 0 changed, 1 destroyed.

		Generating the terraform outputs archive...
		* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
		...

		Backuping the terraform output...
		* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
		...

### Output Commands	

To get help type the command `terratemplate output -h`

	NAME:
	   terratemplate output - execute commands on the terraform output archives

	USAGE:
	   terratemplate output command [command options] [arguments...]

	COMMANDS:
		 extract, e  extract a terratemplate output
		 encrypt, e  encrypt a terratemplate output
		 decrypt, e  decrypt a terratemplate output

	OPTIONS:
	   --help, -h  show help

* Extracting output content with command `terratemplate output extract --template-name QA --instance-id 1 -path \tmp\extracts`

* Decrypting output with command `terratemplate output decrypt --template-name QA --instance-id 1 -vault-password "****"`

* Encrypting output with command `terratemplate output encrypt --template-name QA --instance-id 1 -vault-password "****"`

	An encrypted output must be decrypted before using the `encrypt` command.
