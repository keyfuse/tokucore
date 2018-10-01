// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"fmt"
	"net"
	"time"

	"github.com/tokublock/tokucore/network"
)

func newMockNode(netwrk *network.Network, endpoint string) (*Connection, func()) {
	// Server.
	listener, err := net.Listen("tcp", endpoint)
	if err != nil {
		panic(err)
	}
	listen := func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			c := NewConnection(netwrk, conn, mockServerOnReceive, mockOnDisconnect)
			go c.ReadLoop()
		}
	}
	go listen()

	// Client.
	conn, err := net.DialTimeout("tcp", endpoint, 10*time.Second)
	if err != nil {
		panic(err)
	}
	client := NewConnection(netwrk, conn, mockClientOnReceive, mockOnDisconnect)
	go client.ReadLoop()

	return client, func() {
		client.Close()
		listener.Close()
	}
}

func mockServerOnReceive(c *Connection, msg Message) {
	switch msg.(type) {
	case *MsgVersion:
		rsp := NewMsgVersion(network.TestNet)
		c.Send(rsp)
	case *MsgVerAck:
		rsp := NewMsgVerAck()
		c.Send(rsp)
	}
}

func mockClientOnReceive(c *Connection, msg Message) {
	fmt.Printf("--client.recv:%s, %+v\n", msg.Command(), msg)
}

func mockOnDisconnect(c *Connection) {
	fmt.Printf("[%v]-->[%v].closed\n", c.LocalAddress(), c.RemoteAddress())
}
