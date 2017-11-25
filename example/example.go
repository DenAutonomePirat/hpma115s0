package main

import (
	"encoding/json"
	"fmt"
	"github.com/denautonomepirat/hpma115s0"
	"github.com/jacobsa/go-serial/serial"
	"log"
	"os"
)

func main() {
	fmt.Println("morning")

	options := serial.OpenOptions{
		PortName:        "/dev/tty31",
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	// Make sure to close it later.
	defer port.Close()

	hpma := hpma115s0.NewHpma11520(port)
	m := hpma.ReadParticleMeasurement()
	b, _ := json.Marshal(m)
	os.Stdout.Write(b)
	fmt.Println(m.Pm25)
}
