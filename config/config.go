package config

import (
	"log/slog"
	"reflect"

	"github.com/meian/atgo/logs"
	"github.com/spf13/viper"
)

type config struct {
	Workspace       string `mapstructure:"workspace" default:"."`
	DefaultLogLevel string `mapstructure:"default_log_level" default:"none"`
}

var Config config

func Inititalize() {
	setEnvDefaults(&Config)
	viper.SetEnvPrefix("ATGO")
	viper.AutomaticEnv()
	if err := viper.Unmarshal(&Config); err != nil {
		slog.Warn("failed to read config from environment: %v", err)
		return
	}
	_, err := logs.ParseLevel(Config.DefaultLogLevel)
	if err != nil {
		Config.DefaultLogLevel = "none"
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
