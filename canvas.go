package ruc

import (
	"github.com/fogleman/gg"
	"image/color"
)

type Canvas struct {
	*gg.Context
	mem []uint32
}

func NewCanvas(width, height uint) *Canvas {
	return &Canvas{
		gg.NewContext(int(width), int(height)),
		make([]uint32, width*height),
	}
}

func (c *Canvas) GenSetPixelsMessage(force bool) *SetPixelsMessage {
	msg := NewSetPixelsMessage()
	for y := 0; y < c.Height(); y++ {
		for x := 0; x < c.Width(); x++ {
			if c.pixelUnchanged(x, y, c.Image().At(x, y)) && !force {
				continue
			}
			r, g, b, _ := c.Image().At(x, y).RGBA()
			msg.AppendPixel(Pixel{
				X: uint8(x), Y: uint8(y),
				R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8),
			})
			c.updateMirror(x, y, c.Image().At(x, y))
		}
	}
	return msg
}

func (c *Canvas) comparableColor(color color.Color) uint32 {
	r, g, b, _ := color.RGBA()
	return uint32(uint8(r>>8))<<16 | uint32(uint8(g>>8))<<8 | uint32(uint8(b>>8))
}

func (c *Canvas) pixelUnchanged(x, y int, color color.Color) bool {
	return c.comparableColor(color) == c.mem[y*c.Width()+x]
}

func (c *Canvas) updateMirror(x, y int, color color.Color) {
	c.mem[y*c.Width()+x] = c.comparableColor(color)
}
