package util

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetEnv(name string, path string, stru interface{}) {
	viper.SetConfigName(name)   // name of config file (without extension)
	viper.SetConfigType("yml")  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path)   // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = viper.Unmarshal(&stru)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
