package db

import (
	"testing"
)

func TestRangeReturnsSortedVisibleKeys(t *testing.T) {
	opts := DefaultOptions()
	opts.MemtableMaxRecords = 2
	opts.SSTableCompactThreshold = 10
	db := Open(opts)
	defer db.Close()

	db.Put("b", []byte("value_b"))
	db.Put("d", []byte("value_d"))
	db.Put("a", []byte("value_a"))
	db.Put("c", []byte("value_c"))

	keys := make([]string, 0)
	db.Range("", "", func(key string, value []byte) bool {
		keys = append(keys, key)
		return true
	})

	if len(keys) != 4 {
		t.Errorf("expected 4 keys, got %d", len(keys))
	}
	for i := 1; i < len(keys); i++ {
		if keys[i-1] >= keys[i] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}

	seen := make(map[string]bool)
	for _, key := range keys {
		if seen[key] {
			t.Errorf("duplicate key: %s", key)
		}
		seen[key] = true
	}
}

func TestRangeSkipsTombstonesAndKeepsNewestVersion(t *testing.T) {
	opts := DefaultOptions()
	opts.MemtableMaxRecords = 1
	opts.SSTableCompactThreshold = 10
	db := Open(opts)
	defer db.Close()

	db.Put("x", []byte("old_value"))
	db.Put("y", []byte("value_y"))
	db.Put("z", []byte("old_z"))
	db.Delete("x")
	db.Put("z", []byte("new_z"))

	results := make(map[string]string)
	db.Range("", "", func(key string, value []byte) bool {
		results[key] = string(value)
		return true
	})
	if _, exists := results["x"]; exists {
		t.Error("tombstoned key 'x' should not appear in range")
	}
	if val, ok := results["z"]; !ok || val != "new_z" {
		t.Errorf("expected 'new_z' for key 'z', got '%s'", val)
	}
}

func TestCompactionReplacesTwoOldestSSTables(t *testing.T) {
	opts := DefaultOptions()
	opts.MemtableMaxRecords = 1
	opts.SSTableCompactThreshold = 2
	db := Open(opts)
	defer db.Close()

	db.Put("a", []byte("1"))
	db.Put("b", []byte("2"))
	db.Put("c", []byte("3"))
	db.Put("a", []byte("4"))
	db.Put("b", []byte("5"))

	db.mu.RLock()
	sstableCount := len(db.sstables)
	db.mu.RUnlock()
	if sstableCount > 2 {
		t.Errorf("expected compaction to reduce sstable count, got %d", sstableCount)
	}

	val, ok := db.Get("a")
	if !ok || string(val) != "4" {
		t.Errorf("expected '4' for key 'a', got '%s'", string(val))
	}
	val, ok = db.Get("b")
	if !ok || string(val) != "5" {
		t.Errorf("expected '5' for key 'b', got '%s'", string(val))
	}
	val, ok = db.Get("c")
	if !ok || string(val) != "3" {
		t.Errorf("expected '3' for key 'c', got '%s'", string(val))
	}
}

func TestRangeStopsWhenCallbackReturnsFalse(t *testing.T) {
	opts := DefaultOptions()
	db := Open(opts)
	defer db.Close()

	keys := []string{"a", "b", "c", "d", "e"}
	for _, k := range keys {
		db.Put(k, []byte("value"))
	}

	called := 0
	db.Range("", "", func(key string, value []byte) bool {
		called++
		return called < 3
	})

	if called != 3 {
		t.Errorf("callback called %d times, expected 3", called)
	}
}

func TestCloseStopsFlusherAndCompactor(t *testing.T) {
	opts := DefaultOptions()
	opts.MemtableMaxRecords = 2
	opts.SSTableCompactThreshold = 2
	db := Open(opts)

	for i := 0; i < 100; i++ {
		db.Put(string(rune('a'+i%26)), []byte("value"))
	}
	db.Close()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on Put after Close")
		}
	}()
	db.Put("key", []byte("value"))
}
