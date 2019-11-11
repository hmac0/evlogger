package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
	"github.com/golang/glog"
)

func main() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600, ReadTimeout: time.Millisecond * 500}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	out, err := Execute(s, "sh h")
	if err != nil {
		glog.Error(err)
	}
	
	fmt.Print(out)

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
