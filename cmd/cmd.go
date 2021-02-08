package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/anywhereQL/anywhereQL/cmd/repl"
	"github.com/anywhereQL/anywhereQL/runtime/storage"
	"github.com/anywhereQL/anywhereQL/runtime/storage/aq"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:     "anywhereQL",
		PreRunE: initializeAnywhereQL,
		RunE:    repl.Start,
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "User config (Default: $HOME/.anywhereql.yaml)")

	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Printf("Msg: %s\n", err)
			os.Exit(1)
		}
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/anywhereQL/")
		viper.AddConfigPath(home + "/.anywhereQL")
		viper.AddConfigPath(".")
		viper.SetConfigType("json")
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
		}
	}
}

func initializeAnywhereQL(cmd *cobra.Command, args []string) error {
	se := storage.GetInstance()
	config := viper.AllSettings()
	if _, exists := config["dbs"]; exists {
		for _, v := range config["dbs"].([]interface{}) {
			eng := v.(map[string]interface{})["Engine"].(string)
			sch := v.(map[string]interface{})["Schema"].(string)
			path := v.(map[string]interface{})["Path"].([]interface{})
			pp := []string{}
			for _, p := range path {
				pp = append(pp, p.(string))
			}
			switch eng {
			case "AQDB":
				e, err := aq.Start(pp...)
				if err != nil {
					return err
				}
				se.Add(sch, e)
			}
		}
	}
	return nil
}
