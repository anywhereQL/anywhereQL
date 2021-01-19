package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/anywhereQL/anywhereQL/cmd/repl"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:  "anywhereQL",
		RunE: repl.Start,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "CONFIG_FILE", "User config (Default: $HOME/.anywhereql.yaml)")
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
		viper.AddConfigPath(home)
		viper.SetConfigName(".anywhereql")
	}
	viper.AutomaticEnv()

	viper.ReadInConfig()
}
