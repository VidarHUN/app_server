package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configurations struct {
	Quicrq   QuicrqServer
	Database DatabaseServer
}

type QuicrqServer struct {
	Address string
	Port    int
}

type DatabaseServer struct {
	Address string
	Port    int
}

func ReadConfig(path string) Configurations {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	var configuration Configurations
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	// Set undefined variables
	viper.SetDefault("quicrq.address", "127.0.0.1")
	viper.SetDefault("quicrq.port", 8443)
	viper.SetDefault("database.address", "127.0.0.1")
	viper.SetDefault("database.port", 6379)

	err = viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	fmt.Println("Reading config.yaml file...")
	return configuration
}
