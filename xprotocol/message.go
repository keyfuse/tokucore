// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

// Message --
type Message interface {
	Encode() []byte
	Decode([]byte) error
	Command() string
	Size() int
}

func makeEmptyMessage(command string) Message {
	switch command {
	case CommandVersion:
		return &MsgVersion{}
	case CommandVersionAck:
		return &MsgVerAck{}
	case CommandInventory:
		return &MsgInv{}
	case CommandGetData:
		return &MsgGetData{}
	case CommandGetHeaders:
		return &MsgGetHeaders{}
	case CommandTx:
		return &MsgTx{}
	case CommandHeaders:
		return &MsgHeaders{}
	case CommandReject:
		return &MsgReject{}
	}
	return &MsgUnhandle{command: command}
}
