package kv

import (
    "math/rand"
    "time"
)

const maxLevel = 16
const p = 0.25

type skipListNode struct {
    key, value string
    next       []*skipListNode
}

type SkipListStore struct {
    head  *skipListNode
    level int
}

func NewSkipListStore() *SkipListStore {
    rand.Seed(time.Now().UnixNano())
    return &SkipListStore{
        head:  &skipListNode{next: make([]*skipListNode, maxLevel)},
        level: 1,
    }
}

func randomLevel() int {
    lvl := 1
    for lvl < maxLevel && rand.Float64() < p {
        lvl++
    }
    return lvl
}

func (s *SkipListStore) Set(key, value string) error {
    update := make([]*skipListNode, maxLevel)
    x := s.head
    for i := s.level - 1; i >= 0; i-- {
        for x.next[i] != nil && x.next[i].key < key {
            x = x.next[i]
        }
        update[i] = x
    }
    if x.next[0] != nil && x.next[0].key == key {
        x.next[0].value = value
        return nil
    }
    lvl := randomLevel()
    if lvl > s.level {
        for i := s.level; i < lvl; i++ {
            update[i] = s.head
        }
        s.level = lvl
    }
    newNode := &skipListNode{
        key:   key,
        value: value,
        next:  make([]*skipListNode, lvl),
    }
    for i := 0; i < lvl; i++ {
        newNode.next[i] = update[i].next[i]
        update[i].next[i] = newNode
    }
    return nil
}

func (s *SkipListStore) Get(key string) (string, error) {
    x := s.head
    for i := s.level - 1; i >= 0; i-- {
        for x.next[i] != nil && x.next[i].key < key {
            x = x.next[i]
        }
    }
    x = x.next[0]
    if x != nil && x.key == key {
        return x.value, nil
    }
    return "", ErrNotFound
}

func (s *SkipListStore) Delete(key string) error {
    update := make([]*skipListNode, maxLevel)
    x := s.head
    found := false
    for i := s.level - 1; i >= 0; i-- {
        for x.next[i] != nil && x.next[i].key < key {
            x = x.next[i]
        }
        update[i] = x
    }
    x = x.next[0]
    if x != nil && x.key == key {
        found = true
        for i := 0; i < s.level; i++ {
            if update[i].next[i] != x {
                break
            }
            update[i].next[i] = x.next[i]
        }
        // Decrease level if needed
        for s.level > 1 && s.head.next[s.level-1] == nil {
            s.level--
        }
    }
    if !found {
        return ErrNotFound
    }
    return nil
}

// Range returns all keys in [start, end).
func (s *SkipListStore) Range(start, end string) ([]string, error) {
    var keys []string
    x := s.head
    // Find the first node >= start
    for i := s.level - 1; i >= 0; i-- {
        for x.next[i] != nil && x.next[i].key < start {
            x = x.next[i]
        }
    }
    x = x.next[0]
    for x != nil && x.key < end {
        keys = append(keys, x.key)
        x = x.next[0]
    }
    return keys, nil
}

// Flush is a no-op for in-memory skiplist.
func (s *SkipListStore) Flush() error {
    return nil
}

