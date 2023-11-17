package main

import (
	"github.com/tmsmr/ruc"
	"math/rand"
	"os"
	"time"
)

const TargetFps = 30

type ball struct {
	x, y    float64
	d       float64
	vx, vy  float64
	r, g, b float64
}

func randBall(mx, my int) ball {
	return ball{
		x: float64(mx / 2), y: float64(my / 2),
		d:  float64(rand.Intn((my/2)-2) + 2),
		vx: float64(rand.Intn(mx-10)+10) / float64(mx), vy: float64(rand.Intn(mx-10)+10) / float64(mx),
		r: rand.Float64(), g: rand.Float64(), b: rand.Float64(),
	}
}

var addr string

func init() {
	if len(os.Getenv("UNICORN_ADDR")) != 0 {
		addr = os.Getenv("UNICORN_ADDR")
	} else {
		// emulator
		addr = os.Getenv("localhost:1234")
	}
}

func main() {
	client, err := ruc.NewClient(addr)
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
		start := time.Now()
		c.SetRGB(0, 0, 0)
		c.Clear()
		for i := range balls {
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
		err := client.Display(c, false)
		if err != nil {
			panic(err)
		}
		time.Sleep((time.Second / TargetFps) - time.Since(start))
	}
}
