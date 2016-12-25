package cmd

import (
	"fmt"
	"os"

	"github.com/genzj/go-bingwallpaper/i18n"
	"github.com/genzj/go-bingwallpaper/log"
	"github.com/genzj/go-bingwallpaper/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AppName states name of the project
const AppName string = "go-bingwallpaper"

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
	i18n.LoadTFunc(lang)
	log.Debug(i18n.T("lang_debug_candidate_loaded", i18n.Fields{
		"LangCfg":    lang,
		"LangLoaded": i18n.GetLoadedLang(),
	}))
}

func initConfig() {
	viper.SetEnvPrefix(AppName)
	viper.SetConfigType("json")

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
		log.Fatalf("Specified configuration file %v not readable", cfgFile)

	} else {
		log.Debugln("No config file found, use default settings")
	}
}

func init() {
	cobra.OnInitialize(initI18n)
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config-file", "", "config file (default is ./config.json)")
	RootCmd.PersistentFlags().StringVar(&lang, "lang", "en-us", "language used for display, in xx-YY format.")

}
