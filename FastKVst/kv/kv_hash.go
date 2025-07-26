package kv

import (
    "sync"
)

// HashStore is a simple in‑memory hashmap with RWMutex for concurrency.
type HashStore struct {
    mu    sync.RWMutex
    store map[string]string
}

// NewHashStore constructs a ready‑to‑use HashStore.
func NewHashStore() *HashStore {
    return &HashStore{
        store: make(map[string]string),
    }
}

// Set inserts or updates a key.
func (h *HashStore) Set(key, value string) error {
    h.mu.Lock()
    defer h.mu.Unlock()

    h.store[key] = value
    return nil
}

// Get retrieves a key, or ErrNotFound.
func (h *HashStore) Get(key string) (string, error) {
    h.mu.RLock()
    defer h.mu.RUnlock()

    v, ok := h.store[key]
    if !ok {
        return "", ErrNotFound
    }
    return v, nil
}

// Delete removes a key.
func (h *HashStore) Delete(key string) error {
    h.mu.Lock()
    defer h.mu.Unlock()
	
    if _, ok := h.store[key]; !ok {
        return ErrNotFound
    }
    delete(h.store, key)
    return nil
}

// Range is unsupported for HashStore.
func (h *HashStore) Range(start, end string) ([]string, error) {
    return nil, ErrUnsupported
}

// Flush is a no‑op for in‑memory store; here you could dump to disk.
func (h *HashStore) Flush() error {
    // e.g. write h.store to file if you want
    return nil
}
