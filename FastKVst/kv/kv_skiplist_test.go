package kv

import (
    "sync"
    "testing"
)

func TestSkipListStoreBasic(t *testing.T) {
    s := NewSkipListStore()

    // Set & Get
    if err := s.Set("foo", "bar"); err != nil {
        t.Fatal(err)
    }
    if v, err := s.Get("foo"); err != nil || v != "bar" {
        t.Fatalf("expected bar, got %q, err=%v", v, err)
    }

    // Delete
    if err := s.Delete("foo"); err != nil {
        t.Fatal(err)
    }
    if _, err := s.Get("foo"); err != ErrNotFound {
        t.Fatalf("expected ErrNotFound, got %v", err)
    }

    // Range
    s.Set("a", "1")
    s.Set("b", "2")
    s.Set("c", "3")
    keys, err := s.Range("a", "c")
    if err != nil {
        t.Fatal(err)
    }
    if len(keys) != 2 || keys[0] != "a" || keys[1] != "b" {
        t.Fatalf("unexpected range result: %v", keys)
    }

    // Flush
    if err := s.Flush(); err != nil {
        t.Fatal(err)
    }
}

func TestSkipListStoreConcurrency(t *testing.T) {
    s := NewSkipListStore()
    var wg sync.WaitGroup
    n := 100

    // Writers
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := "key" + string(rune(i))
            value := "val" + string(rune(i))
            if err := s.Set(key, value); err != nil {
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
            _, _ = s.Get(key) // Ignore error, as key may not exist yet
        }(i)
    }

    // Deleters
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := "key" + string(rune(i))
            _ = s.Delete(key) // Ignore error, as key may not exist yet
        }(i)
    }

    wg.Wait()
}