package hashring

import (
	"hash/crc32" // Used to Hash strings to consistent 32 bit Hash
	"log"
	"sort"    // Provides Functions to sort slices and binary Search
	"strconv" // Converts Strings to Integers
)

// HashRing provides a way to initalize a Ring to Add and Remove Nodes
type HashRing struct {
	replicas int // Decided the number of Virtual Nodes per node
	keys     []uint32 // Is a Sorted List of Hash Values (Position on the ring)
	hashMap  map[uint32]string // Maps the Hash Value (pos) to the Node  name
}

// Creates a Hashring with a given number of virtual nodes per Physical node
func New(replicas int) *HashRing {
	return &HashRing{
		replicas: replicas,
		hashMap:  make(map[uint32]string),
	}
}

// Adds one or more node to the ring


func (h *HashRing) Add(nodes ...string) {
	for _, node := range nodes { // For every node
		for i := 0; i < h.replicas; i++ {
			vnode := node + "#" + strconv.Itoa(i) // Create a replica vnode
			sum := h.hash(vnode) // Hash It
			h.keys = append(h.keys, sum) // append the keys
			h.hashMap[sum] = node // Map the Hash to node name
		}
		log.Printf("Added %s with %d replicas", node, h.replicas)
	}
	// Sort the Keys , Sorting is necessary for efficient binary search 
	// when locating the appropriate node for a key.
	sort.Slice(h.keys, func(i, j int) bool { return h.keys[i] < h.keys[j] })
	
}


// Removes a node from the Hash Ring
func (h *HashRing) Remove(node string) {
	for i := 0; i < h.replicas; i++ { // For every replica of this node 
		vnode := node + "#" + strconv.Itoa(i)  // Find the virtual node
		sum := h.hash(vnode) // Hashit
		delete(h.hashMap, sum) // Delete from the HashMap
		// Find the index of the hash in keys
		// Logic below is a bit tricky to understand
		idx := sort.Search(len(h.keys), func(i int) bool { return h.keys[i] >= sum })
		if idx < len(h.keys) && h.keys[idx] == sum {
			h.keys = append(h.keys[:idx], h.keys[idx+1:]...)
		}
	}
	log.Printf("Ring now has %d hash keys", len(h.keys))
}


// Gets the node reponsible for the key
func (h *HashRing) Get(key string) string {
	if len(h.keys) == 0 {
		return ""
	}
	sum := h.hash(key) // Hash the key.
	//Binary search the keys list for the smallest hash ≥ key’s hash.
	idx := sort.Search(len(h.keys), func(i int) bool { return h.keys[i] >= sum })
	//If not found, wrap around to index 0.
	if idx == len(h.keys) {
		idx = 0
	}
	return h.hashMap[h.keys[idx]] // Return the corresponding node from hashMap
}

// Returns a list of all real nodes currently in the hash ring.
func (h *HashRing) Nodes() []string {
	seen := make(map[string]bool)
	result := []string{}
	// Iterate over hashMap.
	for _, v := range h.hashMap {
		//Collect unique real node names.
		if !seen[v] {
			result = append(result, v)
			seen[v] = true
		}
	}
	return result
}

// Hashes a string into a 32-bit integer using CRC32.
func (h *HashRing) hash(data string) uint32 {
	return crc32.ChecksumIEEE([]byte(data))
}
