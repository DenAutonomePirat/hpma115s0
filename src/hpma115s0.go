package hpma115s0

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

const (
	HEAD byte = 0x68

	READ_PARTICLE_MEASURMENT    byte = 0x04
	START_PARTICLE_MEASURMENT   byte = 0x01
	STOP_PARTICLE_MEASURMENT    byte = 0x02
	SET_ADJUSTMENT_COEFFICIENT  byte = 0x08
	READ_ADJUSTMENT_COEFFICIENT byte = 0x08
	STOP_AUTO_SEND              byte = 0x20
	ENABLE_AUTO_SEND            byte = 0x40
)

type Measurement struct {
	TimeStamp int64 `json:"timestamp"`
	Pm25      int   `json:"pm25"`
	Pm10      int   `json:"pm10"`
}

func (m *Measurement) Marshal() *[]byte {
	encoded, _ := json.Marshal(m)
	return &encoded
}

type Hpma115s0 struct {
	Port io.ReadWriteCloser
}

func NewHpma11520(p io.ReadWriteCloser) *Hpma115s0 {
	h := Hpma115s0{
		Port: p,
	}
	return &h
}
func (h *Hpma115s0) name() {

}

func (h *Hpma115s0) ReadParticleMeasurement() *Measurement {
	m := Measurement{
		TimeStamp: int64(time.Now().UnixNano() / 1000 / 1000),
		Pm25:      4,
		Pm10:      1,
	}
	return &m
}

func (h *Hpma115s0) Close() {
	h.Port.Close()
}

func (h *Hpma115s0) SendCmd(c byte, d []byte) {
	cmd := command{
		head: HEAD,
		cmd:  c,
		data: d,
	}
	cmd.len = byte(len(d) + 1)
	fmt.Println(cmd.len)
	cmdb := []byte{cmd.head, cmd.cmd}
	cmdb = append(cmdb, cmd.data...)

	h.Port.Write(cmdb)
}

type command struct {
	head byte
	len  byte
	cmd  byte
	data []byte
	cs   byte
}

func (c *command) check() {

}
