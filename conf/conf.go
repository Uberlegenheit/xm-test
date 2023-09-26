package conf

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		API      API
		LogLevel string
		Postgres Postgres
	}
	API struct {
		ListenOnPort       uint64
		CORSAllowedOrigins []string
	}
	Postgres struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
		SSLMode  string
	}
)

const Service = "xm-task"

func GetNewConfig(path string) (Config, error) {
	// I wasn't sure about "config file" requirement, so I made both .json and .env files
	// I prefer using .env file only, because of its simplicity and convenience
	// Also .env files are cool when using docker-compose
	viper.SetConfigFile(path)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
