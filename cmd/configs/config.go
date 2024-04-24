package configs

import (
	"github.com/spf13/viper"
)

type conf struct {
	RedisAddr            string `mapstructure:"REDIS_ADDR"`
	RedisPassword        string `mapstructure:"REDIS_PASSWORD"`
	RedisDatabase        string `mapstructure:"REDIS_DATABASE"`
	WebServerPort        string `mapstructure:"WEB_SERVER_PORT"`
	LimitRequestPerIP    string `mapstructure:"LIMIT_REQUEST_PER_IP"`
	LimitRequestPerToken string `mapstructure:"LIMIT_REQUEST_PER_TOKEN"`
	BlockingTimeIP       string `mapstructure:"BLOCKING_TIME_IP"`
	BlockingTimeToken    string `mapstructure:"BLOCKING_TIME_TOKEN"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile("./cmd/ratelimitersystem/config.env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, nil
}
