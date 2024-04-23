package config

import (
	"reflect"

	"github.com/spf13/viper"
)

type config struct {
	Workspace string `mapstructure:"workspace" default:"."`
}

var Config config

func Inititalize() {
	setEnvDefaults(&Config)
	viper.SetEnvPrefix("ATGO")
	viper.AutomaticEnv()
	if err := viper.Unmarshal(&Config); err != nil {
		panic(err)
	}
}

func setEnvDefaults(c *config) {
	v := reflect.ValueOf(c).Elem()
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag
		tagValue := tag.Get("default")
		if tagValue != "" {
			viper.SetDefault(tag.Get("mapstructure"), tagValue)
		}
	}
}
