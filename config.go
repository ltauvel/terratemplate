package main

import (
	"encoding/json"
	"io/ioutil"
	"github.com/ltauvel/gomodules/console"
	"github.com/ltauvel/gomodules/filesystem"
)

type TtConfiguration struct {
	Debug string `json:"debug"`
	TemplatesPath string `json:"templates-path"`
	InstancesPath string `json:"instances-path"`
	BackupsPath string `json:"backups-path"`
	VaultPasswordFile string `json:"vault-password-file"`
	ConfigFilePath string `json:"-"`
}

func (configuration *TtConfiguration) Load() {

	// Check if a configuration file exists in the working directory
	cfgpath := filesystem.JoinPath(ttcontext.WorkingDir, "terratemplate.cfg")
	if filesystem.CheckPathExists(cfgpath) {
	
		console.PrintDebug("Loading configuration file '" + cfgpath + "'")
	
		// Reading the configuration file
		file, err := ioutil.ReadFile(cfgpath)

		// Loading JSON content
		err = json.Unmarshal(file, &configuration)
		if err != nil {
			console.PrintError("Error while reading JSON content from configuration file")
		}
	
		// Completing configuration with config file path
		configuration.ConfigFilePath = cfgpath

	}
	
}
