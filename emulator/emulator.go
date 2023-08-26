package main

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type Env struct {
	Debug      bool   `default:"false"`
	BindAddr   string `default:"localhost:1234"`
	MsgTimeout uint   `default:"500"`
	Width      uint   `default:"53"`
	Height     uint   `default:"11"`
	Scaling    uint   `default:"16"`
}

func main() {
	var env Env
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err)
	}
	if env.Debug {
		log.SetLevel(log.DebugLevel)
	}
	display := NewDisplay(env.Width, env.Height, env.Scaling)
	defer display.Close()
	server, err := net.Listen("tcp", env.BindAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("listening for TCP connections on %s", env.BindAddr)
	defer server.Close()
	display.Start()
	for {
		c, err := server.Accept()
		if err != nil {
			log.Warn(err)
			continue
		}
		log.Infof("connection from %s accepted", c.RemoteAddr().String())
		go func(conn net.Conn) {
			defer conn.Close()
			client := NewClient(display.rx)
			client.handle(conn, time.Duration(env.MsgTimeout)*time.Millisecond)
			log.Infof("closing connection from %s", c.RemoteAddr().String())
		}(c)
	}
}
