package main

import (
	"os"

	"github.com/golang/glog"
	"github.com/spf13/viper"
)
func init() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			glog.Error("config not found")
		} else {
			// Config file was found but another error was produced
			glog.Error(err)
		}
		os.Exit(1)
	}

	// TODO: check for required fields 
}
