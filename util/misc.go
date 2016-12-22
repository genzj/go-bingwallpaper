package util

import (
	"github.com/genzj/go-bingwallpaper/log"
	"github.com/kardianos/osext"
)

func ExecutableFolder() string {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}
	return folderPath
}
