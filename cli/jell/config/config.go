package config

import (
	"fmt"
	"github.com/Ygg-Drasill/Jelling/cli/jell/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
)

type JellConfig struct {
	Theme        ui.JellTheme `mapstructure:"theme"`
	SessionToken string       `mapstructure:"sessionToken"`
	BaseApiUrl   string       `mapstructure:"baseUrl"`
}

var DefaultConfig = JellConfig{
	Theme:        ui.Theme,
	SessionToken: "",
	BaseApiUrl:   "http://localhost:30420/api/v1",
}

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	home, err := os.UserHomeDir()
	jellingDirectory := path.Join(home, ".jelling")
	configPath := path.Join(jellingDirectory, "config.yaml")
	cobra.CheckErr(err)

	viper.SetDefault("baseUrl", DefaultConfig.BaseApiUrl)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cobra.CheckErr(os.MkdirAll(jellingDirectory, os.ModePerm))
		viper.Set("theme", ui.Theme)
		cobra.CheckErr(viper.WriteConfigAs(configPath))
	}
	viper.SetConfigFile(configPath)

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
		os.Exit(1)
	}

	conf := JellConfig{}
	err = viper.Unmarshal(&conf)

	ui.Theme = conf.Theme
}
