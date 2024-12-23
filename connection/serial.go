package connection

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

type SerialConnection struct {
	portname string
	speed    int
	port     serial.Port
}

func (s *SerialConnection) SetSpeed(speed string) error {
	var err error
	var temp int
	temp, err = strconv.Atoi(speed)
	if err != nil {
		return errors.New("wrong value for speed")
	}
	s.speed = temp
	return nil
}

func (s *SerialConnection) SetPortName(p string) {
	s.portname = p
}

func (s *SerialConnection) GetPorts() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		return nil, errors.New("no serial ports found")
	} else {
		return ports, nil
	}
}

func (s *SerialConnection) Begin() error {
	var err error
	s.port, err = serial.Open(s.portname, &serial.Mode{})
	if err != nil {
		return err
	}
	mode := &serial.Mode{
		BaudRate: s.speed,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	if err := s.port.SetMode(mode); err != nil {
		return err
	}
	return err
}

func (s *SerialConnection) End() error {
	err := s.port.Close()
	if err != nil {
		return err
	}
	fmt.Printf("Close connection to port %s\n", s.portname)
	return err
}

func (s *SerialConnection) Write(data string) error {
	if !s.IsConnected() {
		return errors.New("not connected")
	}
	log.Debugf("Write data %s", []byte(data))
	_, err := s.port.Write([]byte(data))
	return err
}

func (s *SerialConnection) Read() (string, error) {
	if !s.IsConnected() {
		return "", errors.New("not connected")
	}
	buffer := make([]byte, 255)
	var o string
	t, _ := time.ParseDuration("1ms")
	s.port.SetReadTimeout(t)
	for {
		n, err := s.port.Read(buffer)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 { // EOF
			break
		}
		o = o + string(buffer[:n])
	}
	return o, nil
}

func (s *SerialConnection) IsConnected() bool {
	return s.port != nil
}
