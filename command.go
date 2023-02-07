package waveshare7in5v2

import (
	"github.com/stianeikeland/go-rpio/v4"
)

func (e *Epd) sendCommand(cmd byte) {
	e.dc.Write(rpio.Low)
	e.cs.Write(rpio.Low)

	rpio.SpiTransmit(cmd)

	e.cs.Write(rpio.High)
}

func (e *Epd) sendData(data []byte) {
	e.dc.Write(rpio.High)
	e.cs.Write(rpio.Low)

	for _, chunk := range splitInChunks(data) {
		rpio.SpiTransmit(chunk...)
	}

	e.cs.Write(rpio.High)
}

func (e *Epd) sendCommandWithData(cmd byte, data []byte) {
	e.sendCommand(cmd)
	e.sendData(data)
}

func (e *Epd) waitUntilIdle() {
	for {
		e.sendCommand(GET_STATUS)

		if e.busy.Read() == rpio.High {
			wait(200)

			break
		}

		wait(1)
	}
}
