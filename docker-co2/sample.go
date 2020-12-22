package main

import (
	"log"
	"time"

	"github.com/tarm/serial"
)

func main() {
	var n int
	m := 0
	var err error
	output := []byte{}
	c := &serial.Config{Name: "/dev/ttyS0", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	tmp := []byte{0xff, 0x87, 0x87, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf2}
	n, err = s.Write(tmp)
	if err != nil {
		log.Fatal(err)
	}

	m = 0
	go func() {
		for {
			buf := make([]byte, 128)
			n, err = s.Read(buf)
			if err != nil {
				log.Fatal(err)
				break
			}
			for _, v := range buf[:n] {
				output = append(output, v)
			}
			m += n
			if m > 8 {
				break
			}
		}
	}()
	for {
		if len(output) > 8 {
			break
		}
		time.Sleep(1000)
	}
	log.Printf("%v:%q", m, output[:m])
	time.Sleep(10)

	tmp2 := []byte{0xff, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79}
	n, err = s.Write(tmp2)
	if err != nil {
		log.Fatal(err)
	}
	output = []byte{}
	m = 0
	go func() {
		for {
			buf := make([]byte, 128)
			n, err = s.Read(buf)
			if err != nil {
				log.Fatal(err)
				break
			}
			for _, v := range buf[:n] {
				output = append(output, v)
			}
			m += n
			if m > 8 {
				break
			}
		}
	}()
	for {
		if len(output) > 8 {
			break
		}
		time.Sleep(1000)
	}
	log.Printf("%v:%q", m, output[:m])

}
