package hashring

import (
	"testing"
)
func TestAddAndGet(t *testing.T) {
	ring := New(3)
	ring.Add("NodeA", "NodeB")

	if node := ring.Get("key123"); node == "" {
		t.Error("Expected a node, got empty string")
	}
}

func TestRemoveNode(t *testing.T) {
	ring := New(3)
	ring.Add("NodeA")
	ring.Remove("NodeA")

	if node := ring.Get("key123"); node != "" {
		t.Error("Expected no node after removal, got one")
	}
}

func TestOnlyOneNodeLeft(t *testing.T) {
	ring := New(3)
	ring.Add("NodeA", "NodeB")
	ring.Remove("NodeA")

	nodes := ring.Nodes()
	if len(nodes) != 1 || nodes[0] != "NodeB" {
		t.Errorf("Expected only NodeB, got %v", nodes)
	}
}
