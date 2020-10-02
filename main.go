package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
)

const (
	deviceTypeBMS = "bms"
)

func main() {
	var interval, baud int
	var stopChan = make(chan os.Signal, 2)
	var version, command, serialPath, deviceType, format string
	var versionFlag bool

	// flags
	flag.BoolVar(&versionFlag, "version", false, "prints commit hash for current build")
	flag.IntVar(&baud, "baud", 19200, "baud rate of the device")
	flag.StringVar(&command, "cmd", "help", "talk serial")
	flag.StringVar(&serialPath, "serial-device", "", "the name of the tty device you want to talk to")
	flag.StringVar(&deviceType, "device-type", "bms", "the type of serial device")
	flag.StringVar(&format, "format", "", "the output format for logged data")
	flag.IntVar(&interval, "interval", 1, "time in seconds between logged messages")
	flag.Parse()

	// monitor for signals
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// TODO: set up glog with logDir value from the config

	if versionFlag {
		fmt.Printf("%s", version)
		os.Exit(0)
	}

	if serialPath != "" {

		connect, err := NewConnection(serialPath, baud)
		if err != nil {
			glog.Error(err)
			return
		}

		// encode command as byte array
		bCMD := []byte(command)

		go connect.Start()

		go func(i int, c []byte) {
			// issue command first time
			connect.Write(bCMD)
			for {
				// wait interval
				time.Sleep(time.Duration(interval) * time.Second)

				// issue command
				connect.Write(bCMD)
			}
		}(interval, bCMD)

	}
	<-stopChan // wait for SIGINT
}

// 	if *Cmd != "" {

// 		c := &serial.Config{Name: *SerialDevice, Baud: viper.GetInt("serialBaud"), ReadTimeout: time.Millisecond * 50}
// 		s, err := serial.OpenPort(c)
// 		if err != nil {
// 			glog.Fatal(err)
// 		}

// 		out, err := Execute(s, *Cmd)
// 		if err != nil {
// 			glog.Error(err)
// 		}

// 		glog.V(2).Infof("%s", out)

// 		// paresedOutput := Parse(out, `(c\d{1,2}).{3}(\d.\d{3})`, "bms,cell=%s volts=%s\n")
// 		// if err != nil {
// 		// 	glog.Error(err)

// 		// }
// 		// glog.V(2).Infof("%s", paresedOutput)
// 	}
// }

// // Parse
// func Parse(data string, regex regexp.Regexp, outputPattern string) (string, err) {
// 	//re := regexp.MustCompile(regex)

// 	// `(c\d{1,2}).{3}(\d.\d{3})`\
// 	r := re.FindAllString(input, -1)

// 	glog.V(2).Infof("%v", r)

// 	parsedOutput := ""
// 	for _, elem := range r {
// 		parsedOutput += fmt.Sprintf(outputPattern, elem[1:])
// 	}
// 	return "", parsedOutput

// }

// // Execute takes a pointer to a serial.Port and command string to execute, returns string and error
// func Execute(s *serial.Port, cmd string) (string, error) {
// 	start := time.Now()
// 	_, err := s.Write([]byte(fmt.Sprintf("%s\n", cmd)))
// 	if err != nil {
// 		glog.Error(err)
// 		return "", err
// 	}

// 	reader := bufio.NewReader(s)
// 	buf, err := reader.ReadBytes('\x3E')
// 	if err != nil {
// 		panic(err)
// 	}
// 	//fmt.Println(reply)

// 	// buf := make([]byte, 1024)
// 	// output := ""

// 	// for {
// 	// 	n, _ = s.Read(buf)
// 	// 	if n == 0 {
// 	// 		break
// 	// 	}
// 	// 	glog.Infof("%+v", buf[:n])
// 	// 	output = output + string(buf[:n])
// 	// }

// 	glog.V(2).Infof("execute took %d ms", (time.Since(start)).Nanoseconds()/1000000)

// 	// NOTE: output is trimmed to exclude command and serial cli prompt output
// 	return string(buf[len(cmd):]), nil
// }
