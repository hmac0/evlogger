package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
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

	// TODO: set up glog with logDir value from the config

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
		}

		glog.V(2).Infof("%s", out)

		// paresedOutput := Parse(out, `(c\d{1,2}).{3}(\d.\d{3})`, "bms,cell=%s volts=%s\n")
		// if err != nil {
		// 	glog.Error(err)

		// }
		// glog.V(2).Infof("%s", paresedOutput)
	}
}

// Parse
func Parse(input string, regex string, outputPattern string) string {
	re := regexp.MustCompile(regex)

	// `(c\d{1,2}).{3}(\d.\d{3})`\
	r := re.FindAllString(input, -1)

	glog.V(2).Infof("%v", r)

	parsedOutput := ""
	for _, elem := range r {
		parsedOutput += fmt.Sprintf(outputPattern, elem[1:])
	}
	return parsedOutput

}

// Execute takes a pointer to a serial.Port and command string to execute, returns string and error
func Execute(s *serial.Port, cmd string) (string, error) {
	start := time.Now()
	_, err := s.Write([]byte(fmt.Sprintf("%s\n", cmd)))
	if err != nil {
		glog.Error(err)
		return "", err
	}

	reader := bufio.NewReader(s)
	buf, err := reader.ReadBytes('\x3E')
	if err != nil {
		panic(err)
	}
	//fmt.Println(reply)

	// buf := make([]byte, 1024)
	// output := ""

	// for {
	// 	n, _ = s.Read(buf)
	// 	if n == 0 {
	// 		break
	// 	}
	// 	glog.Infof("%+v", buf[:n])
	// 	output = output + string(buf[:n])
	// }

	glog.V(2).Infof("execute took %d ms", (time.Since(start)).Nanoseconds()/1000000)

	// NOTE: output is trimmed to exclude command and serial cli prompt output
	return string(buf[len(cmd):]), nil
}
