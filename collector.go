package main

import (
	"bufio"
	"fmt"
	"regexp"
	"time"

	"github.com/golang/glog"
	"github.com/tarm/serial"
)

// Collector contains command string, pattern for matching the output of the command, array of ints representing
// the order of submatches in output
type Collector struct {
	SerialConfig       *serial.Config
	SerialPort         *serial.Port
	Command            string
	Pattern            regexp.Regexp
	SubmatchIndexOrder []int
	Output             string
	IntervalMS         int
	Done               chan bool
}

// Start initates the process that collects and logs data
func (c *Collector) Start() {
	ticker := time.NewTicker(time.Duration(c.IntervalMS) * time.Millisecond)
	c.Done = make(chan bool)

	// open serial port
	s, err := serial.OpenPort(c.SerialConfig)
	if err != nil {
		glog.Fatal(err)
	}
	c.SerialPort = s

	go func() {
		for {
			select {
			case <-c.Done:
				return
			case t := <-ticker.C:

				// send command to serial device
				out, err := c.execute()
				if err != nil {
					glog.Error(err)
					return
				}

				// format output from command
				fOut, err := c.format(out)
				if err != nil {
					glog.Error(err)
					return
				}

				// format and save
				//TODO: record fout to a file

			}
		}
	}()
}

// execute sends a command to a serial device and returns the output
func (c *Collector) execute() (string, error) {

	start := time.Now()

	// send command to serial deivice
	_, err := c.SerialPort.Write([]byte(fmt.Sprintf("%s\n", c.Command)))
	if err != nil {
		glog.Error(err)
		return "", err
	}

	// read response from device until EOL
	reader := bufio.NewReader(c.SerialPort)
	buf, err := reader.ReadBytes('\x3E')
	if err != nil {
		panic(err)
	}

	glog.V(2).Infof("execute took %d ms", (time.Since(start)).Nanoseconds()/1000000)

	// NOTE: output is trimmed to exclude command and serial cli prompt output
	return string(buf[len(c.Command):]), nil
}

// parse applys regex and constructs a new string uses submatches and c.output
func (c *Collector) format(data string) (string, error) {

	start := time.Now()

	// apply regex
	r := c.Pattern.FindStringSubmatch(data)
	if r == nil {
		glog.Errorf("no match found")
		return "", fmt.Errorf("no match found")
	}

	glog.V(2).Infof("%v", r)

	// TODO: build a interface of submatches using c.SubmatchIndexOrder

	glog.V(2).Infof("format took %d ms", (time.Since(start)).Nanoseconds()/1000000)

	// TODO: replace with formated output using c.Output and product of above todo
	return "", nil
}
