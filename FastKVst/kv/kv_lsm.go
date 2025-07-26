package kv

import (
    "sort"
    "sync"
)

// lsmEntry represents a key-value pair.
type lsmEntry struct {
    key, value string
}

// LSMStore is a minimal in-memory LSM-tree with two levels.
type LSMStore struct {
    mu        sync.RWMutex
    memtable  map[string]string      // mutable
    immuTable []lsmEntry            // immutable, sorted
    threshold int                   // flush threshold
}

// NewLSMStore constructs a ready-to-use LSMStore.
func NewLSMStore() *LSMStore {
    return &LSMStore{
        memtable:  make(map[string]string),
        threshold: 1000, // flush after 1000 keys (tune as needed)
    }
}

// Set inserts or updates a key.
func (l *LSMStore) Set(key, value string) error {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.memtable[key] = value
    if len(l.memtable) >= l.threshold {
        l.flush()
    }
    return nil
}

// Get retrieves a key, or ErrNotFound.
func (l *LSMStore) Get(key string) (string, error) {
    l.mu.RLock()
    defer l.mu.RUnlock()
    if v, ok := l.memtable[key]; ok {
        return v, nil
    }
    // Binary search in immutable table
    i := sort.Search(len(l.immuTable), func(i int) bool {
        return l.immuTable[i].key >= key
    })
    if i < len(l.immuTable) && l.immuTable[i].key == key {
        return l.immuTable[i].value, nil
    }
    return "", ErrNotFound
}

// Delete marks a key as deleted (tombstone).
func (l *LSMStore) Delete(key string) error {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.memtable[key] = "" // empty string as tombstone
    return nil
}

// Range returns all keys in [start, end).
func (l *LSMStore) Range(start, end string) ([]string, error) {
    l.mu.RLock()
    defer l.mu.RUnlock()
    keySet := make(map[string]struct{})
    // From memtable
    for k := range l.memtable {
        if k >= start && k < end && l.memtable[k] != "" {
            keySet[k] = struct{}{}
        }
    }
    // From immutable table
    for _, e := range l.immuTable {
        if e.key >= start && e.key < end && e.value != "" {
            keySet[e.key] = struct{}{}
        }
    }
    // Collect and sort
    keys := make([]string, 0, len(keySet))
    for k := range keySet {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    return keys, nil
}

// Flush moves memtable to immutable table.
func (l *LSMStore) Flush() error {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.flush()
    return nil
}

// flush (internal): moves memtable to immutable, sorted.
func (l *LSMStore) flush() {
    entries := make([]lsmEntry, 0, len(l.memtable))
    for k, v := range l.memtable {
        entries = append(entries, lsmEntry{k, v})
    }
    sort.Slice(entries, func(i, j int) bool {
        return entries[i].key < entries[j].key
    })
    l.immuTable = mergeLSMEntries(l.immuTable, entries)
    l.memtable = make(map[string]string)
}

// mergeLSMEntries merges two sorted slices, newer values overwrite older.
func mergeLSMEntries(a, b []lsmEntry) []lsmEntry {
    out := make([]lsmEntry, 0, len(a)+len(b))
    i, j := 0, 0
    for i < len(a) && j < len(b) {
        if a[i].key < b[j].key {
            out = append(out, a[i])
            i++
        } else if a[i].key > b[j].key {
            out = append(out, b[j])
            j++
        } else {
            out = append(out, b[j]) // prefer newer
            i++
            j++
        }
    }
    for i < len(a) {
        out = append(out, a[i])
        i++
    }
    for j < len(b) {
        out = append(out, b[j])
        j++
    }
    return out
}