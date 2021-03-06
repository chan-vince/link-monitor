package cmd

import (
	log "github.com/sirupsen/logrus"
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
		Username        string `mapstructure:"username"`
		Password        string `mapstructure:"password"`
		ExchangeName    string `mapstructure:"exchange_name"`
		ExchangeType    string `mapstructure:"exchange_type"`
		RoutingKey      string `mapstructure:"routing_key"`
		PublishInterval uint    `mapstructure:"publish_interval"`
	}
}

func GetConf() *Config {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/")
	viper.SetConfigName("link-monitor.conf")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()

	if err != nil {
		log.Errorf("%v", err)
	}

	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		log.Fatalf("unable to decode into config struct, %v", err)
	}

	return conf
}
