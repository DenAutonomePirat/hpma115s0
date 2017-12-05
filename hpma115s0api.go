package main

import (
	"encoding/json"
	"fmt"
	"github.com/denautonomepirat/hpma115s0/src"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/jacobsa/go-serial/serial"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Println("morning")
	config := client.UDPConfig{
		Addr: "10.0.0.1:8089",
	}
	// Make client
	c, err := client.NewUDPClient(config)
	if err != nil {
		panic(err.Error())
	}

	options := serial.OpenOptions{
		PortName:              "/dev/ttyS0",
		BaudRate:              9600,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       32,
		InterCharacterTimeout: 200,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	// Make sure to close it later.
	defer port.Close()

	hpma := hpma115s0.NewHpma11520(port)

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "us",
	})

	count := 0

	for {

		buf := make([]byte, 40)
		n, err := port.Read(buf)
		if err != nil {
			fmt.Println(err)
		}

		m := hpma.ReadParticleMeasurement()

		m.Pm25 = int((buf[6] * 255) + buf[7])
		m.Pm10 = int((buf[8] * 255) + buf[9])
		b, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%X\n", buf[0:n])

		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))

		// Create a point and add to batch
		tags := map[string]string{"enmvironment": "air"}
		fields := map[string]interface{}{
			"pm25": m.Pm25,
			"pm10": m.Pm10,
		}

		pt, err := client.NewPoint("house", tags, fields, time.Now())
		if err != nil {
			panic(err.Error())
		}
		count++
		bp.AddPoint(pt)
		if count > 10 {
			fmt.Println("influx")
			// Write the batch
			c.Write(bp)
			count = 0
		}

	}

	hpma.Close()
}
