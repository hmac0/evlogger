package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/golang/glog"
	"github.com/tarm/serial"
	//"github.com/golang/regex"
)

// Version holds the version and build commit hash, set at build time.
var Version string

func main() {
	VersionFlag := flag.Bool("version", false, "prints commit hash for current build")
	Cmd := flag.String("cmd", "help", "talk serial")
	flag.Parse()

	// TODO: set up glog with logDir value from the config

	if *VersionFlag {
		fmt.Printf("%s", Version)
		os.Exit(0)
	}

	if *Cmd != "" {

		c := &serial.Config{Name: viper.GetString("serialEVCCPath"), Baud: viper.GetInt("serialBaud"), ReadTimeout: time.Millisecond * 500}
		s, err := serial.OpenPort(c)
		if err != nil {
			glog.Fatal(err)
		}

		for i := 0; i < 5; i++ {

			out, err := Execute(s, *Cmd)
			if err != nil {
				glog.Error(err)
			}

			fmt.Printf("%s", out)
		}
	}
}

func (o *string) Parse() (string, error) {
	reCells := regex.MustCompile(`(c\d{1,2}).{3}(\d.\d{3})`)

}

// Execute takes a pointer to a serial.Port and command string to execute, returns string and error
func Execute(s *serial.Port, cmd string) (string, error) {
	n, err := s.Write([]byte(fmt.Sprintf("%s\n", cmd)))
	if err != nil {
		glog.Error(err)
		return "", err
	}

	buf := make([]byte, 4096)
	output := ""

	for {
		// HACK: used magic number to detect end of output ("evcc> \n" > 8 bytes?)
		n, err = s.Read(buf)
		if n == 0 || n == 8 {
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
