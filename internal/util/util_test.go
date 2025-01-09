package util

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateHash(t *testing.T) {
	uuid1 := GenerateUUIDForClient("192.168.7.1")
	uuid2 := GenerateUUIDForClient("192.168.7.1")
	uuid3 := GenerateUUIDForClient("192.168.7.2")
	// Assert that identical receipts produce the same UUID
	if uuid1 != uuid2 {
		t.Error("UUID for identical receipts should match")
	}

	// Assert that different receipts produce different UUIDs
	if uuid1 == uuid3 {
		t.Error("UUID for different receipts should not match")
	}
}

func TestMemoizationPerformance(t *testing.T) {
	// Helper to generate a list of unique IDs
	generateUniqueIDs := func(count int) []string {
		ids := make([]string, count)
		for i := 0; i < count; i++ {
			ids[i] = uuid.NewString()
		}
		return ids
	}

	// Benchmark for unique IDs
	t.Run("UniqueIDs", func(t *testing.T) {
		uniqueIDs := generateUniqueIDs(1000)

		start := time.Now()
		for _, id := range uniqueIDs {
			_ = GenerateUUIDForClient(id)
		}
		duration := time.Since(start)

		t.Logf("Time taken for 100 unique IDs: %v", duration)
	})

	// Benchmark for repeated IDs
	t.Run("RepeatedIDs", func(t *testing.T) {
		repeatedID := "repeated-id"

		start := time.Now()
		for i := 0; i < 1000; i++ {
			_ = GenerateUUIDForClient(repeatedID)
		}
		duration := time.Since(start)

		t.Logf("Time taken for 100 repeated IDs: %v", duration)
	})
}
