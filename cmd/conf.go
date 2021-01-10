package cmd

import (
	"fmt"
	"github.com/spf13/viper"
)

// Have to use mapstructure when unmarshalling into embedded structs because it doesn't support
// snake case otherwise(!) https://github.com/spf13/viper/issues/125
type Config struct {
	Logging struct {
		Level string `mapstructure:"level"`
		File  string `mapstructure:"file"`
	}

	CacheDir string   `mapstructure:"cache_directory"`
	KitId    string   `mapstructure:"kit_id"`
	Links    []string `mapstructure:"links"`

	Broker   struct {
		Host            string `mapstructure:"host"`
		Port            int    `mapstructure:"port"`
		RoutingKey      string `mapstructure:"routing_key"`
		PublishInterval int    `mapstructure:"publish_interval"`
	}
}

func GetConf() *Config {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/")
	viper.SetConfigName("link-monitor.conf")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("%v", err)
	}

	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}
