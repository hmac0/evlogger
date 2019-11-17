package main

import (
	"github.com/spf13/viper"
	"flag"
	"fmt"
	"os"

	"github.com/tarm/serial"
	"github.com/golang/glog"
)

// Version holds the version and build commit hash, set at build time.
var Version string

func main() {
	VersionFlag := flag.Bool("version", false, "prints commit hash for current build")
	flag.Parse()

	// TODO: set up glog with logDir value from the config 

	if *VersionFlag {
		fmt.Printf("%s\n", Version)
		os.Exit(0)
	}

	c := &serial.Config{Name: viper.GetString("serialBMSPath"), Baud: viper.GetInt("serialBaud"), ReadTimeout: time.Millisecond * 500}
	s, err := serial.OpenPort(c)
	if err != nil {
		glog.Fatal(err)
	}

	out, err := Execute(s, "sh h")
	if err != nil {
		glog.Error(err)
	}
}

// Execute takes a pointer to a serial.Port and command string to execute, returns string and error
func Execute(s *serial.Port, cmd string) (string, error) {
	n, err := s.Write([]byte(fmt.Sprintf("%s\n", cmd)))
	if err != nil {
		glog.Error(err)
		return "", err
	}

	buf := make([]byte, 128)
	output := ""

	for {
		n, err = s.Read(buf)
		if n == 0 {
			break
		}

		if err != nil {
			glog.Error(err)
			return "", err
		}

		output = output + string(buf[:n])
	}

	return output, nil

}
