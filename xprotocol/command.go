// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

const (
	// CommandVersion is the command used when a node creates an outgoing connection.
	// it will immediately advertise its version.
	// The remote node will respond with its version.
	CommandVersion = "version"

	// CommandVersionAck message is sent in reply to version.
	CommandVersionAck = "verack"

	// CommandInventory Allows a node to advertise its knowledge of one or more objects
	// It can be received unsolicited, or in reply to getblocks.
	CommandInventory = "inv"

	// CommandPing is sent primarily to confirm that the TCP/IP connection is still valid.
	CommandPing = "ping"

	// CommandPong  is sent in response to a CommandPing.
	CommandPong = "pong"

	// CommandBlock is sent in response to a getdata message which requests transaction information from a block hash
	CommandBlock = "block"

	// CommandGetBlocks returns an inv packet containing the list of blocks starting right after the
	// last known hash in the block locator object, up to stop value or 500 blocks (max).
	CommandGetBlocks = "getblocks"

	// CommandGetData is used in response to inv, to retrieve the content of a specific object (usually sent after receiving an inv packet).
	CommandGetData = "getdata"

	// CommandGetHeaders return a headers packet containing the headers of blocks starting right after the
	// last known hash in the block locator object, up to stop value or 2000 blocks (max).
	CommandGetHeaders = "getheaders"

	// CommandHeaders returns block headers in response to a getheaders packet
	CommandHeaders = "headers"

	// CommandTx describes a bitcoin transaction, in response to CommandGetData.
	CommandTx = "tx"

	// CommandReject is sent when messages are rejected.
	CommandReject = "reject"
)
