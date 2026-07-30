package db

import (
	"multithreading_hw_3/memtable"
	"sort"
)

func (db *DB) Range(start, end string, fn func(key string, value []byte) bool) {
	db.mu.RLock()
	if db.closed {
		db.mu.RUnlock()
		panic("db closed")
	}
	mutable := db.mutable
	immutables := make([]*memtable.Memtable, len(db.immutables))
	copy(immutables, db.immutables)
	sstables := make([]*SSTable, len(db.sstables))
	copy(sstables, db.sstables)
	db.mu.RUnlock()
	if start != "" && end != "" && start >= end {
		return
	}
	latest := make(map[string]memtable.Record)
	mutable.Range(func(rec memtable.Record) bool {
		if _, exists := latest[rec.Key]; !exists {
			latest[rec.Key] = rec
		}
		return true
	})
	for i := len(immutables) - 1; i >= 0; i-- {
		immutables[i].Range(func(rec memtable.Record) bool {
			if _, exists := latest[rec.Key]; !exists {
				latest[rec.Key] = rec
			}
			return true
		})
	}
	for i := len(sstables) - 1; i >= 0; i-- {
		for _, rec := range sstables[i].records {
			if _, exists := latest[rec.Key]; !exists {
				latest[rec.Key] = rec
			}
		}
	}
	keys := make([]string, 0, len(latest))
	for key, rec := range latest {
		if rec.Tombstone {
			continue
		}
		if start != "" && key < start {
			continue
		}
		if end != "" && key >= end {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		rec := latest[key]
		valueCopy := make([]byte, len(rec.Value))
		copy(valueCopy, rec.Value)
		if !fn(key, valueCopy) {
			break
		}
	}
}
