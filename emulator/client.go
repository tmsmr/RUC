package main

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/tmsmr/ruc"
	"io"
	"net"
	"os"
	"time"
)

type Client struct {
	pix chan ruc.SetPixelsMessage
}

func NewClient(pix chan ruc.SetPixelsMessage) Client {
	return Client{pix}
}

func (client Client) handle(conn net.Conn, timeout time.Duration) {
	for {
		var err error
		// wait for next header
		conn.SetDeadline(time.Time{})
		pre := make([]byte, ruc.HeaderEncodedLen)
		_, err = conn.Read(pre)
		if err != nil {
			if err == io.EOF {
				log.Debugf("EOF for %s", conn.RemoteAddr())
			} else {
				log.Error(err)
			}
			break
		}
		// try to decode header
		var header ruc.Header
		err = header.DecodeHeader(pre)
		if err != nil {
			log.Warn(err)
			continue
		}
		log.Debugf("receiving data (%d bytes) for message type %x from client id %x", header.DataLen, header.Type, header.SourceId)
		// set deadline for data transmission
		conn.SetDeadline(time.Now().Add(timeout))
		// read data
		data := make([]byte, header.DataLen)
		_, err = conn.Read(data)
		if err != nil {
			if err == io.EOF {
				log.Warnf("EOF for %s during data transmission", conn.RemoteAddr())
			}
			if errors.Is(err, os.ErrDeadlineExceeded) {
				log.Warnf("deadline exceeded for %s during data transmission, closing connection", conn.RemoteAddr())
			} else {
				log.Error(err)
			}
			break
		}
		// handle command
		switch header.Type {
		case ruc.MessageTypeSetPixels:
			var msg ruc.SetPixelsMessage
			err = msg.Decode(append(pre, data...))
			if err != nil {
				log.Warn(err)
				continue
			}
			log.Debugf("decoded %s", msg.String())
			client.pix <- msg
		}
	}
}
