package hashring

import (
	"hash/crc32"
	"log"
	"sort"
	"strconv"
	"sync"
)

// HashRing represents a consistent hashing ring
// with virtual node support and thread safety
// using sync.RWMutex for concurrent access.
type HashRing struct {
	replicas int               // Number of virtual nodes per physical node
	keys     []uint32          // Sorted list of virtual node hashes
	hashMap  map[uint32]string // Mapping from virtual node hash to physical node
	nodeSet  map[string]bool   // Tracks added physical nodes to avoid duplicates
	mu       sync.RWMutex      // Read/Write lock for concurrent safety
}

// New initializes a HashRing with the given number of replicas.
func New(replicas int) *HashRing {
	return &HashRing{
		replicas: replicas,
		hashMap:  make(map[uint32]string),
		nodeSet:  make(map[string]bool),
	}
}

// Add inserts one or more nodes into the hash ring.
// Each node is assigned 'replicas' number of virtual nodes.
func (h *HashRing) Add(nodes ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, node := range nodes {
		if h.nodeSet[node] {
			log.Printf("Node %s already exists, skipping", node)
			continue
		}
		h.nodeSet[node] = true
		for i := 0; i < h.replicas; i++ {
			vnode := node + "#" + strconv.Itoa(i)
			sum := h.hash(vnode)
			h.keys = append(h.keys, sum)
			h.hashMap[sum] = node
		}
		log.Printf("Added node %s with %d replicas", node, h.replicas)
	}

	// Sort the keys to maintain ring order
	sort.Slice(h.keys, func(i, j int) bool { return h.keys[i] < h.keys[j] })
	log.Printf("Ring now has %d hash keys", len(h.keys))
}

// Remove deletes a node and all its virtual nodes from the ring.
func (h *HashRing) Remove(node string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.nodeSet[node] {
		log.Printf("Node %s does not exist, skipping removal", node)
		return
	}
	delete(h.nodeSet, node)
	for i := 0; i < h.replicas; i++ {
		vnode := node + "#" + strconv.Itoa(i)
		sum := h.hash(vnode)
		delete(h.hashMap, sum)
		// Remove from sorted keys slice
		idx := sort.Search(len(h.keys), func(i int) bool { return h.keys[i] >= sum })
		if idx < len(h.keys) && h.keys[idx] == sum {
			h.keys = append(h.keys[:idx], h.keys[idx+1:]...)
		}
	}
	log.Printf("Removed node %s and its replicas", node)
	log.Printf("Ring now has %d hash keys", len(h.keys))
}

// Get returns the closest node in the ring responsible for the given key.
func (h *HashRing) Get(key string) string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.keys) == 0 {
		return ""
	}
	sum := h.hash(key)
	idx := sort.Search(len(h.keys), func(i int) bool { return h.keys[i] >= sum })
	if idx == len(h.keys) {
		idx = 0 // Wrap around to the first node
	}
	return h.hashMap[h.keys[idx]]
}

// Nodes returns a deduplicated list of all physical nodes in the ring.
func (h *HashRing) Nodes() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	seen := make(map[string]bool)
	result := []string{}
	for _, v := range h.hashMap {
		if !seen[v] {
			result = append(result, v)
			seen[v] = true
		}
	}
	return result
}

// Keys returns the sorted list of all virtual node hashes in the ring.
func (h *HashRing) Keys() []uint32 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.keys
}

// HashMap returns the internal mapping from virtual node hash to physical node.
func (h *HashRing) HashMap() map[uint32]string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.hashMap
}

// hash is a utility function to compute a consistent hash value
// using crc32.ChecksumIEEE algorithm.
func (h *HashRing) hash(data string) uint32 {
	return crc32.ChecksumIEEE([]byte(data))
}
