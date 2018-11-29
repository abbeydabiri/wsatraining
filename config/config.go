package config

import (
	"bytes"
	"fmt"

	//needed for sqlite3
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

const (
	keySize   = 32
	nonceSize = 24
)

//Config structure
type Config struct {
	Cookie string

	OS, Path, Address, Proxy string
}

var config Config

//Get ...
func Get() *Config {
	return &config
}

//Init ...
func Init(yamlConfig []byte) {

	viper.SetConfigType("yaml")
	viper.SetDefault("address", "127.0.0.1:8000")

	var err error
	if yamlConfig == nil {
		viper.SetConfigName("config")
		viper.AddConfigPath("./")  // optionally look for config in the working directory
		err = viper.ReadInConfig() // Find and read the config file
	} else {
		err = viper.ReadConfig(bytes.NewBuffer(yamlConfig))
	}

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	config = Config{}
	config.OS = viper.GetString("os")
	config.Path = viper.GetString("path")
	config.Proxy = viper.GetString("proxy")
	config.Address = viper.GetString("address")
}
