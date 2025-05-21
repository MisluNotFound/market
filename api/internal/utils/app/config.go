package app

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Addr               string `mapstructure:"address"`
		BaseIP             string `mapstructure:"base_ip"`
		AccessTokenExpire  int    `mapstructure:"access_token_expire"`
		RefreshTokenExpire int    `mapstructure:"refresh_token_expire"`
	} `mapstructure:"server"`

	Database struct {
		Database string `mapstructure:"database"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Charset  string `mapstructure:"charset"`
	} `mapstructure:"database"`

	OSS struct {
		Type    string `mapstructure:"type"`
		MaxSize int    `mapstructure:"max_size"`
		Root    string `mapstructure:"root"`
	} `mapstructure:"oss"`

	ES struct {
		Addresses []string `mapstructure:"addresses"`
	} `mapstructure:"es"`

	Alipay struct {
		APPID string `mapstructure:"app_id"`
	} `mapstructure:"alipay"`
}

var config *Config

func GetConfig() Config {
	if config == nil {
		Init()
	}

	return *config
}

func Init() {
	config = &Config{}
	configPath := os.Getenv("m_market_config")
	if len(configPath) == 0 {
		configPath = "./config.yaml"
	}
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}
}
