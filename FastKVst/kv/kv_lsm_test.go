package kv

import (
	"fmt"
	"sync"
	"testing"
)

func TestLSMStoreBasic(t *testing.T) {
    lsm := NewLSMStore()

    // Set & Get
    if err := lsm.Set("foo", "bar"); err != nil {
        t.Fatal(err)
    }
    if v, err := lsm.Get("foo"); err != nil || v != "bar" {
		fmt.Print(err)
        t.Fatalf("expected bar, got %q, err=%v", v, err)
    }

    // Delete
    if err := lsm.Delete("foo"); err != nil {
		fmt.Print(err)
        t.Fatal(err)
    }
    if v, err := lsm.Get("foo"); err == ErrNotFound || v != "" {
		fmt.Print(err)
        t.Fatalf("expected ErrNotFound or tombstone, got %q, err=%v", v, err)
    }

    // Range
    lsm.Set("a", "1")
    lsm.Set("b", "2")
    lsm.Set("c", "3")
    keys, err := lsm.Range("a", "c")
    if err != nil {
        t.Fatal(err)
    }
    if len(keys) != 2 || keys[0] != "a" || keys[1] != "b" {
        t.Fatalf("unexpected range result: %v", keys)
    }

    // Flush
    if err := lsm.Flush(); err != nil {
        t.Fatal(err)
    }
}

func TestLSMStoreConcurrency(t *testing.T) {
    lsm := NewLSMStore()
    var wg sync.WaitGroup
    n := 100

    // Writers
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := "key" + string(rune(i))
            value := "val" + string(rune(i))
            if err := lsm.Set(key, value); err != nil {
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
            _, _ = lsm.Get(key) // Ignore error, as key may not exist yet
        }(i)
    }

    // Deleters
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := "key" + string(rune(i))
            _ = lsm.Delete(key) // Ignore error, as key may not exist yet
        }(i)
    }

    wg.Wait()
}