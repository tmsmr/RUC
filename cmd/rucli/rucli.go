package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/tmsmr/ruc"
	"image"
	"os"
	"strconv"
	"strings"
)

type Env struct {
	Address string `default:"localhost:1234"`
	Size    string `default:"53x11"`
}

func set(c *ruc.Client, x, y, r, g, b uint8) error {
	msg := ruc.NewSetPixelsMessage()
	msg.AppendPixel(ruc.Pixel{X: x, Y: y, R: r, G: g, B: b})
	return c.Send(msg)
}

func load(c *ruc.Client, path string, width, height int) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}
	canvas := ruc.NewCanvas(uint(width), uint(height))
	canvas.DrawImage(img, 0, 0)
	return c.Display(canvas, true)
}

func clear(c *ruc.Client, width, height int) error {
	msg := ruc.NewSetPixelsMessage()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			msg.AppendPixel(ruc.Pixel{X: uint8(x), Y: uint8(y)})
		}
	}
	return c.Send(msg)
}

func uint8arg(arg string) uint8 {
	val, err := strconv.Atoi(arg)
	if err != nil {
		panic(err)
	}
	return uint8(val)
}
func main() {
	var env Env
	if err := envconfig.Process("ruc", &env); err != nil {
		panic(err)
	}
	if len(os.Args) < 2 {
		panic("missing command")
	}
	size := strings.Split(env.Size, "x")
	width, err := strconv.Atoi(size[0])
	if err != nil {
		panic(err)
	}
	height, err := strconv.Atoi(size[1])
	if err != nil {
		panic(err)
	}
	c, err := ruc.NewClient(env.Address)
	if err != nil {
		panic(err)
	}
	if err = c.Connect(); err != nil {
		panic(err)
	}
	defer c.Close()
	switch os.Args[1] {
	case "set":
		if len(os.Args) < 7 {
			panic("missing arguments")
		}
		if err = set(c, uint8arg(os.Args[2]), uint8arg(os.Args[3]), uint8arg(os.Args[4]), uint8arg(os.Args[5]), uint8arg(os.Args[6])); err != nil {
			panic(err)
		}
		break
	case "load":
		if len(os.Args) < 3 {
			panic("missing arguments")
		}
		if err = load(c, os.Args[2], width, height); err != nil {
			panic(err)
		}
		break
	case "clear":
		if err = clear(c, width, height); err != nil {
			panic(err)
		}
		break
	default:
		panic("unknown command")
	}
}
