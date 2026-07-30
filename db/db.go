package db

import (
	"sync"

	"multithreading_hw_3/memtable"
)

type DB struct {
	mu   sync.RWMutex
	opts Options

	mutable    *memtable.Memtable
	immutables []*memtable.Memtable
	sstables   []*SSTable

	flushCh   chan *memtable.Memtable
	compactCh chan struct{}
	closeCh   chan struct{}
	wg        sync.WaitGroup

	closed bool
}

func Open(opts Options) *DB {
	if opts.MemtableMaxRecords <= 0 {
		opts = DefaultOptions()
	}
	db := &DB{
		opts:    opts,
		mutable: memtable.New(),

		immutables: make([]*memtable.Memtable, 0),
		sstables:   make([]*SSTable, 0),

		flushCh:   make(chan *memtable.Memtable, 100),
		compactCh: make(chan struct{}),
		closeCh:   make(chan struct{}),

		closed: false,
	}
	db.wg.Add(2)
	go db.runFlusher()
	go db.runCompactor()
	return db
}

func (db *DB) Close() {
	db.mu.Lock()
	if db.closed {
		db.mu.Unlock()
		return
	}
	db.closed = true
	if db.mutable.Len() > 0 {
		immutable := db.mutable
		db.immutables = append(db.immutables, immutable)
		db.mutable = nil
		db.mu.Unlock()
		db.flushCh <- immutable
	} else {
		db.mutable = nil
		db.mu.Unlock()
	}
	close(db.closeCh)
	db.wg.Wait()
}

func (db *DB) Put(key string, value []byte) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.closed {
		panic("db closed")
	}
	copyValue := make([]byte, len(value))
	copy(copyValue, value)
	db.mutable.Put(key, copyValue)
	if db.mutable.Len() >= db.opts.MemtableMaxRecords {
		immutable := db.mutable
		db.immutables = append(db.immutables, immutable)
		db.mutable = memtable.New()
		db.flushCh <- immutable
	}
}

func (db *DB) Get(key string) ([]byte, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.closed {
		panic("db closed")
	}
	if isTombstoneInMemtable(db.mutable, key) {
		return nil, false
	}
	if val, ok := db.mutable.Get(key); ok {
		copyVal := make([]byte, len(val))
		copy(copyVal, val)
		return copyVal, true
	}
	for i := len(db.immutables) - 1; i >= 0; i-- {
		if isTombstoneInMemtable(db.immutables[i], key) {
			return nil, false
		}
		if val, ok := db.immutables[i].Get(key); ok {
			copyVal := make([]byte, len(val))
			copy(copyVal, val)
			return copyVal, true
		}
	}
	for i := len(db.sstables) - 1; i >= 0; i-- {
		if rec, ok := db.sstables[i].Get(key); ok {
			if rec.Tombstone {
				return nil, false
			}
			copyVal := make([]byte, len(rec.Value))
			copy(copyVal, rec.Value)
			return copyVal, true
		}
	}
	return nil, false
}

func (db *DB) Delete(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.closed {
		panic("DB IS closed")
	}
	db.mutable.Delete(key)
	if db.mutable.Len() >= db.opts.MemtableMaxRecords {
		immutable := db.mutable
		db.immutables = append(db.immutables, immutable)
		db.mutable = memtable.New()
		db.flushCh <- immutable
	}
}

func isTombstoneInMemtable(mt *memtable.Memtable, key string) bool {
	found := false
	isTombstone := false
	mt.Range(func(rec memtable.Record) bool {
		if rec.Key == key {
			found = true
			isTombstone = rec.Tombstone
			return false
		}
		if rec.Key > key {
			return false
		}
		return true
	})
	return found && isTombstone
}
