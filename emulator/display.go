package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	log "github.com/sirupsen/logrus"
	"github.com/tmsmr/ruc"
	"os"
)

type Display struct {
	rx chan ruc.SetPixelsMessage
	w  int
	h  int
	s  int
	fb []byte
}

func NewDisplay(w uint, h uint, s uint) Display {
	fb := make([]byte, w*h*4)
	d := Display{
		make(chan ruc.SetPixelsMessage),
		int(w), int(h), int(s), fb,
	}
	ebiten.SetWindowTitle("RUC Emulator")
	ebiten.SetWindowSize(d.w*int(s), d.h*int(s))
	ebiten.SetWindowClosingHandled(true)
	return d
}

func (d Display) Start() {
	log.Infof("starting display (%dx%d, scale %d)", d.w, d.h, d.s)
	go d.handleRx()
	go d.gameLoop()
}

func (d Display) handleRx() {
	for {
		msg := <-d.rx
		log.Debugf("updating %d pixels", len(msg.Pixels))
		for _, p := range msg.Pixels {
			if int(p.X) >= d.w || int(p.Y) >= d.h {
				log.Warnf("dropping invalid pixel %s", p.String())
				continue
			}
			pos := (int(p.Y)*d.w + int(p.X)) * 4
			copy(d.fb[pos:pos+3], p.RGB())
		}
	}
}

func (d Display) gameLoop() {
	if err := ebiten.RunGame(&d); err != nil {
		log.Fatal(err)
	}
}

func (d Display) Close() {
	close(d.rx)
}

func (d Display) Update() error {
	if ebiten.IsWindowBeingClosed() {
		log.Info("closing window")
		os.Exit(0)
	}
	return nil
}

func (d Display) Draw(canvas *ebiten.Image) {
	canvas.WritePixels(d.fb)
}

func (d Display) Layout(_, _ int) (int, int) {
	return d.w, d.h
}
