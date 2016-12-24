package cmd

import (
	"github.com/genzj/go-bingwallpaper/i18n"
	"github.com/genzj/go-bingwallpaper/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(coreCmd)
}

var coreCmd = &cobra.Command{
	Use:   "core",
	Short: "Trigger pic downloading and other core features",
	Long:  `Trigger pic downloading and other core features`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"test": viper.GetString("test"),
		}).Infoln(i18n.T("hello_world"))
		log.Infoln("Bye.")
		execute(cmd, args)
	},
}

func execute(cmd *cobra.Command, args []string) {
	log.WithFields(log.Fields{
		"test": viper.GetString("test"),
	}).Infoln("Hello world!")
	log.Infoln("Bye.")
}
