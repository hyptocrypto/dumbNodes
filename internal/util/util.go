package util

import (
	"crypto/sha1"
	"sync"

	"github.com/google/uuid"
)

const RecNameSpace = "123e4567-e89b-12d3-a456-426614174000"

var (
	ClientIdMemo = map[string]uuid.UUID{}
	ClientMemMux = sync.RWMutex{}
)

// GenerateUUIDForClient creates a deterministic UUID for a given id.
func GenerateUUIDForClient(id string) uuid.UUID {
	ClientMemMux.RLock()
	if id, ok := ClientIdMemo[id]; ok {
		ClientMemMux.RUnlock()
		return id
	}
	ClientMemMux.RUnlock()
	var data []byte
	data = append(data, []byte(id)...)
	namespace := uuid.MustParse(RecNameSpace)
	newId := uuid.NewHash(sha1.New(), namespace, data, 5)
	ClientMemMux.Lock()
	ClientIdMemo[id] = newId
	ClientMemMux.Unlock()
	return newId
}
