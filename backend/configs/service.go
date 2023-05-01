package configs

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("configs")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Config file: %s not found", "path"))
		} else {
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
	}
}
