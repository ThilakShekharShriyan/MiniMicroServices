package kv

import "testing"

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
