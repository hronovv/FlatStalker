package cache

import (
	"testing"
	"time"
)

func TestLRUGetPutEvict(t *testing.T) {
	c := NewLRU[string](2, time.Minute)

	c.Put(1, "a")
	c.Put(2, "b")
	if v, ok := c.Get(1); !ok || v != "a" {
		t.Fatalf("get 1: %q %v", v, ok)
	}

	c.Put(3, "c") // evicts 2 (least recently used after Get(1))
	if _, ok := c.Get(2); ok {
		t.Fatal("expected key 2 evicted")
	}
	if v, ok := c.Get(3); !ok || v != "c" {
		t.Fatalf("get 3: %q %v", v, ok)
	}
}

func TestLRUTTL(t *testing.T) {
	c := NewLRU[string](8, 20*time.Millisecond)
	c.Put(1, "a")
	if _, ok := c.Get(1); !ok {
		t.Fatal("expected hit")
	}
	time.Sleep(30 * time.Millisecond)
	if _, ok := c.Get(1); ok {
		t.Fatal("expected expired miss")
	}
}
