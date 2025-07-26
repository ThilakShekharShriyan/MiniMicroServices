package kv

import (
    "sync"
)

// trieNode represents a node in the Trie.
type trieNode struct {
    children map[rune]*trieNode
    value    *string
}

// TrieStore is a thread-safe Trie-based KV store.
type TrieStore struct {
    mu   sync.RWMutex
    root *trieNode
}

// NewTrieStore constructs a ready-to-use TrieStore.
func NewTrieStore() *TrieStore {
    return &TrieStore{
        root: &trieNode{children: make(map[rune]*trieNode)},
    }
}

// Set inserts or updates a key.
func (t *TrieStore) Set(key, value string) error {
    t.mu.Lock()
    defer t.mu.Unlock()
    node := t.root
    for _, ch := range key {
        if node.children[ch] == nil {
            node.children[ch] = &trieNode{children: make(map[rune]*trieNode)}
        }
        node = node.children[ch]
    }
    node.value = &value
    return nil
}

// Get retrieves a key, or ErrNotFound.
func (t *TrieStore) Get(key string) (string, error) {
    t.mu.RLock()
    defer t.mu.RUnlock()
    node := t.root
    for _, ch := range key {
        if node.children[ch] == nil {
            return "", ErrNotFound
        }
        node = node.children[ch]
    }
    if node.value == nil {
        return "", ErrNotFound
    }
    return *node.value, nil
}

// Delete removes a key.
func (t *TrieStore) Delete(key string) error {
    t.mu.Lock()
    defer t.mu.Unlock()
    var parents []*trieNode
    node := t.root
    for _, ch := range key {
        if node.children[ch] == nil {
            return ErrNotFound
        }
        parents = append(parents, node)
        node = node.children[ch]
    }
    if node.value == nil {
        return ErrNotFound
    }
    node.value = nil
    // Optional: prune empty nodes (not implemented for simplicity)
    return nil
}

// Range returns all keys in [start, end).
func (t *TrieStore) Range(start, end string) ([]string, error) {
    t.mu.RLock()
    defer t.mu.RUnlock()
    var keys []string
    var dfs func(node *trieNode, prefix string)
    dfs = func(node *trieNode, prefix string) {
        if node.value != nil && prefix >= start && prefix < end {
            keys = append(keys, prefix)
        }
        for ch, child := range node.children {
            dfs(child, prefix+string(ch))
        }
    }
    dfs(t.root, "")
    return keys, nil
}

// Flush is a no-op for in-memory Trie.
func (t *TrieStore) Flush() error {
    return nil
}