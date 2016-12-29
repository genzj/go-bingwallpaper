package cmd

import (
	"fmt"
	"os"

	"github.com/genzj/gobingwallpaper/i18n"
	"github.com/genzj/gobingwallpaper/log"
	"github.com/genzj/gobingwallpaper/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AppName states name of the project
const AppName string = "gobingwallpaper"

var cfgFile string
var lang string

// RootCmd is the entry of whole program
var RootCmd = &cobra.Command{
	Use:   AppName,
	Short: "A Go application downloading latest wallpaper from Bing.com.",
	Long: `Download the daily bing wallpaper and set it your desktop. Also
integrates a WEB UI for configuration.`,
}

//Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func initI18n() {
	i18n.SetLanguageFilePath(util.ExecutableFolder() + "/i18n")
	i18n.LoadTFunc(viper.GetString("global.Language"))
	log.Debug(i18n.T("lang_debug_candidate_loaded", i18n.Fields{
		"LangCfg":    viper.GetString("global.Language"),
		"LangLoaded": i18n.GetLoadedLang(),
	}))
}

func initConfig() {
	viper.SetEnvPrefix(AppName)
	viper.SetConfigType("yaml")

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")                // name of config file (without extension)
		viper.AddConfigPath(util.ExecutableFolder()) // adding home directory as first search path
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debugln("Using config file:", viper.ConfigFileUsed())
	} else if cfgFile != "" {
		log.WithFields(log.Fields{
			"CfgFile": cfgFile,
		}).Fatalf(i18n.T("config_file_not_readable"))

	} else {
		log.Debugln("No config file found, use default settings")
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initI18n)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config-file", "", "config file (default is ./config.yaml under app installation path.)")
	RootCmd.PersistentFlags().String("history-file", "", "history file (default is ./history.json under app installation path.)")
	viper.BindPFlag("global.HistoryFile", RootCmd.PersistentFlags().Lookup("history-file"))
	viper.SetDefault("global.HistoryFile", "./history.json")
	RootCmd.PersistentFlags().StringVar(&lang, "lang", "en-us", "language used for display, in xx-YY format.")
	viper.BindPFlag("global.Language", RootCmd.PersistentFlags().Lookup("lang"))
	viper.SetDefault("global.Language", "en-us")
}
