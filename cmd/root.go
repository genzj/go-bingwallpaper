package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   APP_NAME,
	Short: "A Go application downloading latest wallpaper from Bing.com.",
	Long: `Download the daily bing wallpaper and set it your desktop. Also
integrates a WEB UI for configuration.`,
}
