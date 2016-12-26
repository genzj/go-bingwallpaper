package util

import (
	"github.com/genzj/gobingwallpaper/log"
	"github.com/kardianos/osext"
)

// ExecutableFolder returns path to the folder containing currently running
// executable file
func ExecutableFolder() string {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}
	return folderPath
}
