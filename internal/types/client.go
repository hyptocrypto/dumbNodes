package types

import (
	"github.com/google/uuid"
)

// Id will be a hash of the ip addr, and key will be some encryption key.
// TODO: Some key crypto bruuu
type ClientConn struct {
	ClientId uuid.UUID
	Key      string
}

type ClientConnections map[uuid.UUID]*ClientConn
