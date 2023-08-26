package ruc

import (
	"net"
)

type Client struct {
	id    byte
	raddr *net.TCPAddr
	conn  net.Conn
}

func NewClient(unicorn string) (*Client, error) {
	return NewClientWithId(unicorn, 0)
}

func NewClientWithId(unicorn string, id byte) (*Client, error) {
	raddr, err := net.ResolveTCPAddr("tcp", unicorn)
	if err != nil {
		return nil, err
	}
	return &Client{
		id:    id,
		raddr: raddr,
	}, nil
}

func (c *Client) Connect() error {
	conn, err := net.DialTCP("tcp", nil, c.raddr)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Send(msg Message) error {
	msg.SetSourceId(c.id)
	data := msg.Encode()
	_, err := c.conn.Write(data)
	return err
}
func (c *Client) Display(canvas *Canvas, force bool) error {
	msg := canvas.GenSetPixelsMessage(force)
	return c.Send(msg)
}

func (c *Client) Close() error {
	return c.conn.Close()
}
