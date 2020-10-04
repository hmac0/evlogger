package main

import (
	"bytes"
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

		connect, err := NewConnection(serialPath, baud, 32)
		if err != nil {
			glog.Error(err)
			return
		}

		// encode command as byte array
		bCMD := []byte(command + "\n")

		go connect.Start()

		go func(i int, c []byte) {
			for {
				// wait interval
				time.Sleep(time.Duration(interval) * time.Second)

				// dump lineChan
				// for line := range connect.lineChan {
				// 	glog.V(2).Infof("dropping: %s", line)
				// }

				done := false
				bStop := []byte{98, 109, 115, 99, 62, 32}

				// issue command
				connect.Write(bCMD)

				// keep going until done or timeout
				for !done {

					glog.V(2).Infof("buffer length: %d", len(connect.lineChan))

					select {
					case line := <-connect.lineChan:
						if bytes.Equal(line, bStop) {
							glog.V(2).Info("done")
							done = true
							break
						}

						glog.Infof("%s", string(line))
						glog.V(2).Infof("%v", line)

					case <-time.After(time.Duration(2*interval) * time.Second):
						glog.Warning("timed out!")
						done = true
					}
				}
			}
		}(interval, bCMD)

	}
	<-stopChan // wait for SIGINT
}