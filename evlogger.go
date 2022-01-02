package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/golang/glog"
	"github.com/tarm/serial"
)

// Version holds the version and build commit hash, set at build time.
var Version string

func main() {
	VersionFlag := flag.Bool("version", false, "prints commit hash for current build")
	Cmd := flag.String("cmd", "help", "talk serial")
	SerialDevice := flag.String("serial_device", "/dev/ttyUSB0", "the name of the tty device you want to talk to")
	flag.Parse()

	if *VersionFlag {
		fmt.Printf("%s", Version)
		os.Exit(0)
	}

	if *Cmd != "" {
		c := &serial.Config{Name: *SerialDevice, Baud: viper.GetInt("serialBaud"), ReadTimeout: time.Millisecond * 50}
		s, err := serial.OpenPort(c)
		if err != nil {
			glog.Fatal(err)
		}

		out, err := Execute(s, *Cmd)
		if err != nil {
			glog.Error(err)
		} else {
			glog.Infof("%s", out)
		}
	}
}

// Execute takes a pointer to a serial.Port and command string to execute, returns string and error
func Execute(s *serial.Port, cmd string) (string, error) {

	// record start time for performance monitoring
	start := time.Now()

	// write command as a byte array into serial port
	_, err := s.Write([]byte(fmt.Sprintf("%s\n", cmd)))
	if err != nil {
		glog.Error(err)
		return "", err
	}

	// TODO: should the buffer be created ahead of sending the command
	// TODO: break up sending the command and receiving the output

	// init new buffer for serial output
	// NOTE: '\x3E' is a magic char, that should indicate the end of the output
	reader := bufio.NewReader(s)
	buf, err := reader.ReadBytes('\x3E')
	if err != nil {
		glog.Error(err)
	}

	glog.V(2).Infof("execute took %d ms", (time.Since(start)).Nanoseconds()/1000000)

	// NOTE: output is trimmed to exclude command and serial cli prompt output
	return string(buf[len(cmd):]), nil
}
