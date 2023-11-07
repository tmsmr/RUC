package ruc

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	HeaderPreambleVal = 0xAA
	HeaderEncodedLen  = 10
	PixelEncodedLen   = 5
)

const (
	MessageTypeSetPixels byte = iota
)

var (
	ErrEncHeaderLen     = errors.New("invalid length for encoded header")
	ErrChecksumMismatch = errors.New("invalid checksum")
	ErrEncMessageLen    = errors.New("invalid length for encoded message")
	ErrEncPixelLen      = errors.New("invalid length for encoded pixel")
)

type Message interface {
	GetType() byte
	SetSourceId(id byte)
	GetSourceId() byte
	Encode() []byte
	Decode(data []byte) error
}

type Header struct {
	SourceId byte
	Type     byte
	DataLen  uint16
}

func (h *Header) EncodeHeader() []byte {
	data := []byte{
		HeaderPreambleVal,
		HeaderPreambleVal,
		HeaderPreambleVal,
		h.SourceId,
		0x00,
		h.Type,
		0x00,
	}
	dl := make([]byte, 2)
	binary.BigEndian.PutUint16(dl, h.DataLen)
	data = append(data, dl...)
	var checksum byte
	for _, v := range data {
		checksum ^= v
	}
	return append(data, checksum)
}

func (h *Header) DecodeHeader(data []byte) error {
	if len(data) != HeaderEncodedLen {
		return ErrEncHeaderLen
	}
	var checksum byte
	for i := 0; i < HeaderEncodedLen-1; i++ {
		checksum ^= data[i]
	}
	if checksum != data[len(data)-1] {
		return ErrChecksumMismatch
	}
	h.SourceId = data[3]
	h.Type = data[5]
	h.DataLen = binary.BigEndian.Uint16(data[7:9])
	return nil
}

type Pixel struct {
	X uint8
	Y uint8
	R uint8
	G uint8
	B uint8
}

func (p *Pixel) String() string {
	return fmt.Sprintf("%d,%d (%d,%d,%d)", p.X, p.Y, p.R, p.G, p.B)
}

func (p *Pixel) RGB() []byte {
	return []byte{p.R, p.G, p.B}
}

func (p *Pixel) Encode() []byte {
	return []byte{p.X, p.Y, p.R, p.G, p.B}
}

func (p *Pixel) Decode(data []byte) error {
	if len(data) != PixelEncodedLen {
		return ErrEncPixelLen
	}
	p.X = data[0]
	p.Y = data[1]
	p.R = data[2]
	p.G = data[3]
	p.B = data[4]
	return nil
}

type SetPixelsMessage struct {
	Header
	Pixels []Pixel
}

func NewSetPixelsMessage() *SetPixelsMessage {
	return &SetPixelsMessage{
		Header: Header{
			Type:    MessageTypeSetPixels,
			DataLen: 0,
		},
		Pixels: make([]Pixel, 0),
	}
}

func (m *SetPixelsMessage) String() string {
	return fmt.Sprintf("SetPixelsMessage (%d pixels)", len(m.Pixels))
}

func (m *SetPixelsMessage) GetType() byte {
	return MessageTypeSetPixels
}

func (m *SetPixelsMessage) SetSourceId(id byte) {
	m.SourceId = id
}

func (m *SetPixelsMessage) GetSourceId() byte {
	return m.SourceId
}

func (m *SetPixelsMessage) Encode() []byte {
	data := m.EncodeHeader()
	for _, p := range m.Pixels {
		data = append(data, p.Encode()...)
	}
	return data
}

func (m *SetPixelsMessage) Decode(data []byte) error {
	if len(data) < HeaderEncodedLen {
		return ErrEncHeaderLen
	}
	m.DecodeHeader(data[0:HeaderEncodedLen])
	dl := m.DataLen
	for pos := 0; pos < int(dl)/PixelEncodedLen; pos++ {
		pos := HeaderEncodedLen + pos*PixelEncodedLen
		if len(data) < pos+PixelEncodedLen {
			return ErrEncMessageLen
		}
		var p Pixel
		p.Decode(data[pos : pos+PixelEncodedLen])
		m.AppendPixel(p)
	}
	return nil
}

func (m *SetPixelsMessage) AppendPixel(p Pixel) {
	m.Pixels = append(m.Pixels, p)
	m.DataLen = uint16(len(m.Pixels) * PixelEncodedLen)
}
