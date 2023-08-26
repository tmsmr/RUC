package main

import (
	"github.com/tmsmr/ruc"
	"math/rand"
	"time"
)

type ball struct {
	x, y    float64
	d       float64
	vx, vy  float64
	r, g, b float64
}

func randBall(mx, my int) ball {
	return ball{
		x: rand.Float64() * float64(mx), y: rand.Float64() * float64(my),
		d:  rand.Float64() * 8,
		vx: rand.Float64(), vy: rand.Float64(),
		r: rand.Float64(), g: rand.Float64(), b: rand.Float64(),
	}
}

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
	c := ruc.NewCanvas(53, 11)

	balls := []ball{
		randBall(c.Width(), c.Height()),
		randBall(c.Width(), c.Height()),
		randBall(c.Width(), c.Height()),
	}

	c.SetRGB(0, 0, 0)
	c.Clear()
	// force transmitting all pixels, since we don't know the current state of the display
	if err = client.Display(c, true); err != nil {
		panic(err)
	}

	for {
		c.SetRGB(0, 0, 0)
		c.Clear()
		for i, _ := range balls {
			// apply vx to x
			if balls[i].x+balls[i].d/2 > float64(c.Width()) || balls[i].x-balls[i].d/2 < 0 {
				balls[i].vx = -balls[i].vx
			}
			balls[i].x += balls[i].vx
			// apply vy to y
			if balls[i].y+balls[i].d/2 > float64(c.Height()) || balls[i].y-balls[i].d/2 < 0 {
				balls[i].vy = -balls[i].vy
			}
			balls[i].y += balls[i].vy
			// draw ball
			c.SetRGB(balls[i].r, balls[i].g, balls[i].b)
			c.DrawCircle(balls[i].x, balls[i].y, balls[i].d/2)
			c.Fill()
		}
		client.Display(c, false)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second / 30)
	}
}
