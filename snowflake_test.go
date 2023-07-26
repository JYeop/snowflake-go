package snowflake

import (
	"testing"
	"time"
)

func TestNode(t *testing.T) {
	_, err := Node(Config{NodeId: 0})
	if err != nil {
		t.Errorf("New returned error: %v", err)
	}
}

func TestGenerate(t *testing.T) {
	sf, err := Node(Config{NodeId: 0})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	id1, err := sf.Generate()
	if err != nil {
		t.Errorf("Generate returned error: %v", err)
	}

	id2, err := sf.Generate()
	if err != nil {
		t.Errorf("Generate returned error: %v", err)
	}

	if id1 == id2 {
		t.Errorf("Generate returned the same ID twice: %v", id1)
	}
}

func TestGetTimeFromId(t *testing.T) {
	sf, err := Node(Config{NodeId: 0})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	id, _ := sf.Generate()
	timestamp := sf.GetTimeFromId(id)
	now := time.Now().UnixMilli()

	// Allow 2ms difference to account for the time it takes to run the code
	if abs(now-timestamp) > 2 {
		t.Errorf("GetTimeFromId returned incorrect time: got %v, want %v", timestamp, now)
	}
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
