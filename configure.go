package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.ReadInConfig()

	env := "dev"
	if envVar := os.Getenv("ENV"); envVar != "" {
		env = strings.ToLower(envVar)
	}

	viper.SetConfigName(env)
	viper.MergeInConfig()
	viper.SetConfigName(fmt.Sprintf("local-%s", env))
	viper.MergeInConfig()

	viper.SetEnvPrefix("oversight")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
