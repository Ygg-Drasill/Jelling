package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Ygg-Drasill/Jelling/cli/jell/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Read or edit jell config",
	Run: func(cmd *cobra.Command, args []string) {
		var conf config.JellConfig
		err := viper.Unmarshal(&conf)
		if err != nil {
			panic(err)
		}
		pretty, err := json.Marshal(conf)
		fmt.Println(string(pretty))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
