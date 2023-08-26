package main

import (
	"github.com/tmsmr/ruc"
	"math/rand"
	"time"
)

func main() {
	// connect to the emulator
	client, err := ruc.NewClient("localhost:1234")
	if err != nil {
		panic(err)
	}
	if err = client.Connect(); err != nil {
		panic(err)
	}
	defer client.Close()

	// galactic unicorn
	bounds := []int{53, 11}
	var x uint8 = 0
	var y uint8 = 0
	for {
		msg := ruc.NewSetPixelsMessage()
		// clear last pixel
		msg.AppendPixel(ruc.Pixel{
			X: x, Y: y,
		})
		// new random position
		x = uint8(rand.Intn(bounds[0]))
		y = uint8(rand.Intn(bounds[1]))
		// add new pixel with random color
		msg.AppendPixel(ruc.Pixel{
			X: x, Y: y,
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
		})
		// transmit changes
		if err = client.Send(msg); err != nil {
			panic(err)
		}
		// wait half a second
		time.Sleep(500 * time.Millisecond)
	}
}
