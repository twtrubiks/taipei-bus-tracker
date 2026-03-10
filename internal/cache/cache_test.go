package cache

import (
	"sync"
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	c := New(1 * time.Second)
	defer c.Close()

	c.Set("key1", "value1")
	v, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected key1 to exist")
	}
	if v != "value1" {
		t.Errorf("expected value1, got %v", v)
	}
}

func TestGetMiss(t *testing.T) {
	c := New(1 * time.Second)
	defer c.Close()

	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("expected miss for nonexistent key")
	}
}

func TestTTLExpiry(t *testing.T) {
	c := New(50 * time.Millisecond)
	defer c.Close()

	c.Set("key1", "value1")

	// Should exist immediately
	if _, ok := c.Get("key1"); !ok {
		t.Fatal("expected key1 to exist before TTL")
	}

	// Wait for expiry
	time.Sleep(100 * time.Millisecond)

	if _, ok := c.Get("key1"); ok {
		t.Error("expected key1 to be expired after TTL")
	}
}

func TestOverwrite(t *testing.T) {
	c := New(1 * time.Second)
	defer c.Close()

	c.Set("key1", "v1")
	c.Set("key1", "v2")

	v, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected key1 to exist")
	}
	if v != "v2" {
		t.Errorf("expected v2, got %v", v)
	}
}

func TestConcurrentAccess(t *testing.T) {
	c := New(1 * time.Second)
	defer c.Close()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			c.Set("key", i)
		}(i)
		go func() {
			defer wg.Done()
			c.Get("key")
		}()
	}
	wg.Wait()
	// No race condition panic = pass
}
