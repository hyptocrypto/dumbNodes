package mproto

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/hyptocrypto/dumbNodes/internal/types"
	"github.com/vmihailenco/msgpack/v5"
)

// Serialize will take the underlying struct instance and create a message pack message based on it.
// The first 4 bytes represent the Response size.
func Serialize(a any) ([]byte, error) {
	p, err := msgpack.Marshal(a)
	if err != nil {
		return nil, err
	}
	pLen := len(p)
	payLoadBuffer := make([]byte, 4+pLen)

	binary.BigEndian.PutUint32(payLoadBuffer[:4], uint32(pLen))
	copy(payLoadBuffer[4:], p)

	return payLoadBuffer, nil
}

func readFromConn(conn net.Conn) ([]byte, error) {
	// Read the length prefix (4 bytes)
	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(conn, lenBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read length prefix: %w", err)
	}

	pLen := binary.BigEndian.Uint32(lenBuf)

	// Read the actual payload
	payload := make([]byte, pLen)
	_, err = io.ReadFull(conn, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to read payload: %w", err)
	}

	return payload, nil
}

func DeserializetReqeustFromConn(conn net.Conn) (*types.Request, error) {
	payload, err := readFromConn(conn)
	if err != nil {
		return nil, err
	}

	var req types.Request
	err = msgpack.Unmarshal(payload, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	return &req, nil
}

func DeserializeResponseFromConn(conn net.Conn) (*types.Response, error) {
	payload, err := readFromConn(conn)
	if err != nil {
		return nil, err
	}

	var resp types.Response
	err = msgpack.Unmarshal(payload, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}
