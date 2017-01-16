package output

import (
	"os"
	"time"
	"github.com/ltauvel/gomodules/console"
	"github.com/ltauvel/gomodules/crypto"
	"github.com/ltauvel/gomodules/targz"
	"github.com/ltauvel/gomodules/filesystem"
	"github.com/ltauvel/terratemplate/secret"
)

type Output struct {
	Handler IOutputHandler
	TemplateName string
	InstanceId string
	Name string
	FullName string
	ModTime time.Time
}

func (output Output) Copy(destination string, force bool) Output { 
	return output.Handler.Copy(output, destination, force)
}

func (output Output) Delete() {
	console.PrintDebug("Deleting output", output.FullName)
	os.Remove(output.FullName)
}

func (output Output) IsEncrypted() bool { 
	console.PrintDebug("Checking that output", output.FullName, "is encrypted")
	return crypto.IsEncrypted(output.FullName)
}

func (output Output) Encrypt(passphrase string) bool { 
	console.PrintDebug("Encrypt output", output.FullName)
	
	if !crypto.IsEncrypted(output.FullName) {
		console.PrintDebug("Processing encryption")
		
		if !crypto.Encrypt(output.FullName, secret.Format(passphrase, secret.HIDE_NONE))  {
			return false
		}
		
		console.PrintDebug("Output successfully encrypted")
	} else {
		console.PrintDebug("Output is already encrypted")
	}
	
	return true
}

func (output Output) Decrypt(passphrase string) bool { 
	console.PrintDebug("Decrypting output", output.FullName)
	
	if crypto.IsEncrypted(output.FullName) {
		console.PrintDebug("Processing decryption")
		
		if !crypto.Decrypt(output.FullName, secret.Format(passphrase, secret.HIDE_NONE))  {
			console.PrintError("The specified vault password is incorrect.")
			
			return false
		}
		
		console.PrintDebug("Output successfully decrypted")
	} else {
		console.PrintDebug("Output is not encrypted")
	}
	
	return true
}

func (output Output) Extract(destination string, passphrase string) bool {
	
	// Creating the destination directory if it does not exists
	filesystem.CreateDir(destination)
	
	// Copying the output archive to the destination directory
	tempoutput := output.Copy(filesystem.JoinPath(destination, output.Name), true)

	// Decrypting the copied output
	tempoutput.Decrypt(secret.Format(passphrase, secret.HIDE_NONE))
	
	// Extracting the copied output
	console.PrintDebug("Extracting output", output.FullName)
	targz.Extract(tempoutput.FullName, destination)
	
	// Removing the trailing copied output
	tempoutput.Delete()
	
	return true
}
