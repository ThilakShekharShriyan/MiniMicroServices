package kv

import "errors"

var (
    ErrNotFound    = errors.New("key not found")
    ErrUnsupported = errors.New("operation unsupported")
)

// KVStore defines the minimal operations for a key-value store.
type KVStore interface {
    // Set writes the key→value pair. Overwrites if key exists.
    Set(key, value string) error

    // Get returns the value for a key, or ErrNotFound.
    Get(key string) (string, error)

    // Delete removes the key from the store.
    Delete(key string) error

    // Range returns all values for keys in [start, end).
    // If unsupported, implement as ErrUnsupported.
    Range(start, end string) ([]string, error)

    // Flush simulates persisting in-memory state to disk.
    Flush() error
}
