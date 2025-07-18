package hashring

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestAddAndGet(t *testing.T) {
	ring := New(3)

	ring.Add("NodeA", "NodeB", "NodeC")

	nodes := ring.Nodes()
	if len(nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(nodes))
	}

	key := "my-key"
	node := ring.Get(key)
	if node == "" {
		t.Errorf("expected a node for key '%s', got empty string", key)
	}
}

func TestRemove(t *testing.T) {
	ring := New(3)

	ring.Add("NodeA", "NodeB")
	ring.Get("hello")

	ring.Remove("NodeA")
	nodes := ring.Nodes()
	if len(nodes) != 1 || nodes[0] != "NodeB" {
		t.Errorf("expected only NodeB to remain, got %+v", nodes)
	}

	// Still get a result even after removing one node
	if ring.Get("hello") == "" {
		t.Errorf("expected a node after removal, got none")
	}

	// Removing again should be a no-op
	ring.Remove("NodeA")
}

func TestDuplicateAdd(t *testing.T) {
	ring := New(3)
	ring.Add("NodeA", "NodeA")

	if len(ring.Nodes()) != 1 {
		t.Errorf("expected only 1 unique node, got %d", len(ring.Nodes()))
	}
}

func TestEmptyRing(t *testing.T) {
	ring := New(3)
	if ring.Get("test") != "" {
		t.Errorf("expected empty response for Get on empty ring")
	}
}
func TestThreadSafety(t *testing.T) {
	ring := New(10)
	var wg sync.WaitGroup

	// Add nodes concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ring.Add("Node" + strconv.Itoa(i))
		}(i)
	}

	// Remove nodes concurrently
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Millisecond * 5) // Let some adds happen first
			ring.Remove("Node" + strconv.Itoa(i))
		}(i)
	}

	// Lookup concurrently
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = ring.Get("key" + strconv.Itoa(i))
		}(i)
	}

	wg.Wait()
}
