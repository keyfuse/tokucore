// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"fmt"
	"net"

	"github.com/tokublock/tokucore/network"
)

// ReceiveFunc --
type ReceiveFunc func(*Connection, Message)

// DisconnectFunc --
type DisconnectFunc func(*Connection)

// Connection --
type Connection struct {
	conn         net.Conn
	stream       *Stream
	onReceive    ReceiveFunc
	onDisconnect DisconnectFunc
}

// NewConnection -- creates new Connection.
func NewConnection(net *network.Network, conn net.Conn, onReceive ReceiveFunc, onDisconnect DisconnectFunc) *Connection {
	return &Connection{
		conn:         conn,
		onReceive:    onReceive,
		onDisconnect: onDisconnect,
		stream:       NewStream(conn, net.Magic),
	}
}

// ReadLoop -- ready for data received.
// This method must in go routine.
func (c *Connection) ReadLoop() {
	for {
		msg, err := c.stream.ReadMessage()
		if err != nil {
			fmt.Printf("connection.ReadLoop.error:%+v\n", err)
			break
		}
		c.onReceive(c, msg)
	}
	c.onDisconnect(c)
}

// RemoteAddress -- remote address.
func (c *Connection) RemoteAddress() string {
	return c.conn.RemoteAddr().String()
}

// LocalAddress -- local address.
func (c *Connection) LocalAddress() string {
	return c.conn.LocalAddr().String()
}

// Send -- send command and datas.
func (c *Connection) Send(msg Message) error {
	return c.stream.WriteMessage(msg)
}

// Close -- close the connection.
func (c *Connection) Close() {
	c.conn.Close()
}
