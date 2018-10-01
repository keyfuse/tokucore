// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xprotocol

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Payload as follows:
//
//        ===============================================
//        |   Field         |   Data Type   |   Size    |
//        ===============================================
//        |   Magic         |   int32       |   4       |
//        -----------------------------------------------
//        |   Command       |   string      |   12      |
//        -----------------------------------------------
//        |   Size          |   int32       |   4       |
//        -----------------------------------------------
//        |   Checksum      |   []byte      |   4       |
//        -----------------------------------------------
//        |   Data          |   []byte      |   ?       |
//        -----------------------------------------------
//

const (
	bufferSize   = 32 * 1024
	readTimeout  = 30 * time.Second
	writeTimeout = 20 * time.Second
)

// Stream --
type Stream struct {
	magic  []byte
	conn   net.Conn
	mutex  sync.Mutex
	reader *bufio.Reader
	writer *bufio.Writer
}

// NewStream -- creates a new stream.
func NewStream(conn net.Conn, magic []byte) *Stream {
	return &Stream{
		conn:   conn,
		magic:  magic,
		reader: bufio.NewReaderSize(conn, bufferSize),
		writer: bufio.NewWriterSize(conn, bufferSize),
	}
}

// ReadMessage -- read a message from the network.
func (s *Stream) ReadMessage() (Message, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.conn.SetReadDeadline(time.Now().Add(readTimeout))

	magic := make([]byte, 4)
	if _, err := io.ReadFull(s.reader, magic); err != nil {
		return nil, err
	}
	command := make([]byte, 12)
	if _, err := io.ReadFull(s.reader, command); err != nil {
		return nil, err
	}
	var size uint32
	if err := binary.Read(s.reader, binary.LittleEndian, &size); err != nil {
		return nil, err
	}
	check := make([]byte, 4)
	if _, err := io.ReadFull(s.reader, check); err != nil {
		return nil, err
	}
	data := make([]byte, size)
	if _, err := io.ReadFull(s.reader, data); err != nil {
		return nil, err
	}
	if !bytes.Equal(check, checksum(data)) {
		return nil, fmt.Errorf("stream.read.payload.checksum.mismatch")
	}

	cmdStr := string(bytes.TrimRight(command, "\x00"))
	msg := makeEmptyMessage(cmdStr)
	if err := msg.Decode(data); err != nil {
		return nil, err
	}
	return msg, nil
}

// WriteMessage -- write the message to the network.
func (s *Stream) WriteMessage(msg Message) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.conn.SetWriteDeadline(time.Now().Add(writeTimeout))

	var b bytes.Buffer
	b.Write(s.magic)
	command := make([]byte, 12)
	copy(command, []byte(msg.Command()))
	b.Write(command)

	datas := msg.Encode()
	binary.Write(&b, binary.LittleEndian, uint32(len(datas)))

	b.Write(checksum(datas))
	b.Write(datas)

	s.writer.Write(b.Bytes())
	return s.writer.Flush()
}

func checksum(data []byte) []byte {
	first := sha256.Sum256(data)
	sum := sha256.Sum256(first[:])
	return sum[0:4]
}
