// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xnet

import (
	"net"
)

// ReceiveFunc --
type ReceiveFunc func(*Connection, Message)

// DisconnectFunc --
type DisconnectFunc func(*Connection)

// Connection --
type Connection struct {
	conn         net.Conn
	network      *Network
	stream       *Stream
	onReceive    ReceiveFunc
	onDisconnect DisconnectFunc
	lastError    error
}

// NewConnection -- creates new Connection.
func NewConnection(network *Network, conn net.Conn, onReceive ReceiveFunc, onDisconnect DisconnectFunc) *Connection {
	return &Connection{
		conn:         conn,
		network:      network,
		onReceive:    onReceive,
		onDisconnect: onDisconnect,
		stream:       NewStream(conn, network.Magic),
	}
}

// ReadLoop -- ready for data received.
// This method must in go routine.
func (c *Connection) ReadLoop() {
	for {
		msg, err := c.stream.ReadMessage()
		if err != nil {
			c.lastError = err
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
