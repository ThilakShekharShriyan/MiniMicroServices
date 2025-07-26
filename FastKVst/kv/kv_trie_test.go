package kv

import (
    "sync"
    "testing"
)

func TestTrieStoreBasic(t *testing.T) {
    trie := NewTrieStore()

    // Set & Get
    if err := trie.Set("foo", "bar"); err != nil {
        t.Fatal(err)
    }
    if v, err := trie.Get("foo"); err != nil || v != "bar" {
        t.Fatalf("expected bar, got %q, err=%v", v, err)
    }

    // Delete
    if err := trie.Delete("foo"); err != nil {
        t.Fatal(err)
    }
    if _, err := trie.Get("foo"); err != ErrNotFound {
        t.Fatalf("expected ErrNotFound, got %v", err)
    }

    // Range
    trie.Set("a", "1")
    trie.Set("b", "2")
    trie.Set("c", "3")
    keys, err := trie.Range("a", "c")
    if err != nil {
        t.Fatal(err)
    }
    if len(keys) != 2 || keys[0] != "a" || keys[1] != "b" {
        t.Fatalf("unexpected range result: %v", keys)
    }

    // Flush
    if err := trie.Flush(); err != nil {
        t.Fatal(err)
    }
}

func TestTrieStoreConcurrency(t *testing.T) {
    trie := NewTrieStore()
    var wg sync.WaitGroup
    n := 100

    // Writers
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := "key" + string(rune(i))
            value := "val" + string(rune(i))
            if err := trie.Set(key, value); err != nil {
                t.Errorf("Set failed: %v", err)
            }
        }(i)
    }

    // Readers
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := "key" + string(rune(i))
            _, _ = trie.Get(key) // Ignore error, as key may not exist yet
        }(i)
    }

    // Deleters
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := "key" + string(rune(i))
            _ = trie.Delete(key) // Ignore error, as key may not exist yet
        }(i)
    }
	wg.Wait()	
}