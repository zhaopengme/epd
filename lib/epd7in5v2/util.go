package epd7in5v2

import (
	"time"

	rpio "github.com/stianeikeland/go-rpio/v4"
)

func digitalWrite(i uint8, state rpio.State) {

	pin := rpio.Pin(i)
	pin.Write(state)
}

func digitalRead(i uint8) rpio.State {
	return rpio.ReadPin(rpio.Pin(i))
}

func delayMS(ms int64) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func spiWrite(command ...byte) {
	err := rpio.SpiBegin(rpio.Spi0)
	if err != nil {
		panic(err)
	}

	rpio.SpiTransmit(command...)
	rpio.SpiEnd(rpio.Spi0)

}
