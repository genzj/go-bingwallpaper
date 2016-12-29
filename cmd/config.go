package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/genzj/gobingwallpaper/log"
	"github.com/genzj/gobingwallpaper/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration manipulation and examination commands",
	Long:  `Configuration manipulation and examination commands`,
	Run:   executeConfig,
}

func executeConfig(cmd *cobra.Command, args []string) {
	if l, _ := cmd.Flags().GetBool("list"); l {
		list()
		os.Exit(0)
	} else if toSet, err := cmd.Flags().GetStringSlice("set"); err == nil {
		set(toSet)
		os.Exit(0)
	}
}

func dumpConfig() ([]byte, error) {
	var conf interface{}
	if err := viper.Unmarshal(&conf); err != nil {
		return []byte{}, err
	}
	if b, err := yaml.Marshal(conf); err != nil {
		return []byte{}, err
	} else {
		return b, nil
	}
}

func list() {
	if b, err := dumpConfig(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%v\n", string(b))
	}
}

func set(toSet []string) {
	parser := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9._]*)=(.+)$`)
	for _, s := range toSet {
		if !parser.MatchString(s) {
			log.WithFields(log.Fields{
				"set": s,
			}).Fatalf("illegal option setting string '%v'", s)
		}
		tokens := parser.FindStringSubmatch(s)
		viper.Set(tokens[1], tokens[2])
		log.Infoln(tokens[1], tokens[2])
	}
	if b, err := dumpConfig(); err != nil {
		log.Fatal(err)
	} else {
		outf := viper.ConfigFileUsed()
		if outf == "" {
			outf = util.ExecutableFolder() + "/config.yaml"
		}
		if info, err := os.Stat(outf); err != nil && !os.IsNotExist(err) {
			log.WithFields(log.Fields{
				"cfgFile": outf,
			}).Fatal(err)
		} else {
			var mode os.FileMode
			if info == nil {
				mode = os.FileMode(0644)
			} else {
				mode = info.Mode()
			}
			if err := ioutil.WriteFile(outf, b, mode); err != nil {
				log.WithFields(log.Fields{
					"cfgFile": outf,
				}).Fatal(err)
			}
		}
	}
}

func init() {
	configCmd.Flags().BoolP("list", "l", false, "show all configuration values then exit")
	configCmd.Flags().StringSlice("set", nil, "set options and save to configuration file, can be repeated several times to set multiple options. configuration must be in path.to.key=value notation, e.g. global.historyfile=his.json")
	RootCmd.AddCommand(configCmd)
}
