package main

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			glog.Error("config not found")
			fmt.Println("config not found")
		} else {
			// Config file was found but another error was produced
			glog.Error(err)
		}
	}
}
