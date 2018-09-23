// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xnet

import (
	"fmt"
	"net"
	"time"
)

func newMockNode(network *Network, endpoint string) (*Connection, func()) {
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
			c := NewConnection(network, conn, mockServerOnReceive, mockOnDisconnect)
			go c.ReadLoop()
		}
	}
	go listen()
	time.Sleep(time.Millisecond * 500)

	// Client.
	conn, err := net.DialTimeout("tcp", endpoint, 10*time.Second)
	if err != nil {
		panic(err)
	}
	client := NewConnection(network, conn, mockClientOnReceive, mockOnDisconnect)
	go client.ReadLoop()
	time.Sleep(time.Millisecond * 500)

	return client, func() {
		client.Close()
		listener.Close()
	}
}

func mockServerOnReceive(c *Connection, msg Message) {
	switch msg.(type) {
	case *MsgVersion:
		rsp := NewMsgVersion(TestNet)
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
