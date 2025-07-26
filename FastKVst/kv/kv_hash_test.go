package kv

import (
	"sync"
	"testing"
)

func TestHashStoreConcurrency(t *testing.T) {
    hs := NewHashStore()
    var wg sync.WaitGroup
    n := 100

    // Writer goroutines
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := "key" + string(rune(i))
            value := "val" + string(rune(i))
            if err := hs.Set(key, value); err != nil {
                t.Errorf("Set failed: %v", err)
            }
        }(i)
    }

    // Reader goroutines
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := "key" + string(rune(i))
            _, _ = hs.Get(key) // Ignore error, as key may not exist yet
        }(i)
    }

    // Deleter goroutines
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := "key" + string(rune(i))
            _ = hs.Delete(key) // Ignore error, as key may not exist yet
        }(i)
    }

    wg.Wait()
}

func TestHashStoreBasic(t *testing.T) {
    hs := NewHashStore()

    // Set & Get
    if err := hs.Set("foo", "bar"); err != nil {
        t.Fatal(err)
    }
    if v, err := hs.Get("foo"); err != nil || v != "bar" {
        t.Fatalf("expected bar, got %q, err=%v", v, err)
    }

    // Delete
    if err := hs.Delete("foo"); err != nil {
        t.Fatal(err)
    }
    if _, err := hs.Get("foo"); err != ErrNotFound {
        t.Fatalf("expected ErrNotFound, got %v", err)
    }

    // Range unsupported
    if _, err := hs.Range("a", "z"); err != ErrUnsupported {
        t.Fatalf("expected ErrUnsupported, got %v", err)
    }
}
