package kv

import (
    "github.com/google/btree"
    "sync"
)

// btreeItem wraps a key-value pair for the B-tree.
type btreeItem struct {
    key, value string
}

func (a btreeItem) Less(b btree.Item) bool {
    return a.key < b.(*btreeItem).key
}

// BTreeStore is a thread-safe B-tree-based KV store.
type BTreeStore struct {
    mu   sync.RWMutex
    tree *btree.BTree
}

// NewBTreeStore constructs a ready-to-use BTreeStore.
func NewBTreeStore() *BTreeStore {
    return &BTreeStore{
        tree: btree.New(2), // degree 2 is minimal; increase for performance
    }
}

func (b *BTreeStore) Set(key, value string) error {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.tree.ReplaceOrInsert(&btreeItem{key, value})
    return nil
}

func (b *BTreeStore) Get(key string) (string, error) {
    b.mu.RLock()
    defer b.mu.RUnlock()
    item := b.tree.Get(&btreeItem{key: key})
    if item == nil {
        return "", ErrNotFound
    }
    return item.(*btreeItem).value, nil
}

func (b *BTreeStore) Delete(key string) error {
    b.mu.Lock()
    defer b.mu.Unlock()
    item := b.tree.Delete(&btreeItem{key: key})
    if item == nil {
        return ErrNotFound
    }
    return nil
}

// Range returns all keys in [start, end).
func (b *BTreeStore) Range(start, end string) ([]string, error) {
    b.mu.RLock()
    defer b.mu.RUnlock()
    var keys []string
    b.tree.AscendRange(&btreeItem{key: start}, &btreeItem{key: end}, func(i btree.Item) bool {
        keys = append(keys, i.(*btreeItem).key)
        return true
    })
    return keys, nil
}

// Flush is a no-op for in-memory B-tree.
func (b *BTreeStore) Flush() error {
    return nil
}