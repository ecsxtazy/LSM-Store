package db

import (
	"sync"

	"multithreading_hw_3/memtable"
)

// DB manages a mutable memtable, a queue of immutable memtables,
// and a list of immutable in-memory SSTables.
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

// Open initializes DB and starts one background flusher and one background compactor.
func Open(opts Options) *DB {
	// TODO:
	// 1) normalize opts using DefaultOptions() when values <= 0.
	// 2) initialize mutable memtable + channels.
	// 3) start runFlusher() and runCompactor() goroutines.
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
	db.wg.Add(1)
	return db
}

func (db *DB) Close() {
	// TODO:
	// 1) set closed flag under lock.
	db.mu.Lock()
	// 2) rotate non-empty mutable into immutables and enqueue for flush.
	if db.closed == true {
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
	// 3) stop background workers and wait for db.wg.
	close(db.closeCh)
	db.wg.Wait()
}

// Put writes a value for key into the current mutable memtable.
func (db *DB) Put(key string, value []byte) {
	// TODO:
	// 1) panic if db is closed.
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.closed {
		panic("db closed")
	}
	// 2) write copy(value) into mutable memtable.
	copyValue := make([]byte, len(value))
	copy(copyValue, value)
	db.mutable.Put(key, copyValue)
	// 3) rotate mutable->immutable after insert when Len reaches MemtableMaxRecords.
	// 4) enqueue immutable into flushCh without blocking forever.
	if db.mutable.Len() >= db.opts.MemtableMaxRecords {
		immutable := db.mutable
		db.immutables = append(db.immutables, immutable)
		db.mutable = memtable.New()
		db.flushCh <- immutable
	}
}

// Get searches data in order: mutable -> immutables -> sstables.
func (db *DB) Get(key string) ([]byte, bool) {
	// TODO:
	// 1) panic if db is closed.
	db.mu.RLock()
	if db.closed {
		panic("db closed")
	}
	// 2) snapshot pointers under lock.
	mutable := db.mutable
	immutables := make([]*memtable.Memtable, len(db.immutables))
	copy(immutables, db.immutables)
	sstables := make([]*SSTable, len(db.sstables))
	copy(sstables, db.sstables)
	db.mu.RUnlock()
	// 3) search from newest to oldest across all layers.
	// 4) return copy(value), tombstone means not found.
	if isTombstoneInMemtable(mutable, key) {
		return nil, false
	}
	if value, ok := mutable.Get(key); ok {
		copyValue := make([]byte, len(value))
		copy(copyValue, value)
		return copyValue, true
	}
	for i := len(immutables) - 1; i >= 0; i-- {
		if isTombstoneInMemtable(immutables[i], key) {
			return nil, false
		}
		if value, ok := immutables[i].Get(key); ok {
			copyValue := make([]byte, len(value))
			copy(copyValue, value)
			return value, true
		}
	}
	for i := len(sstables) - 1; i >= 0; i-- {
		if rec, ok := sstables[i].Get(key); ok {
			if rec.Tombstone {
				return nil, false
			}
			copyValue := make([]byte, len(rec.Value))
			copy(copyValue, rec.Value)
			return rec.Value, true
		}
	}
	return nil, false
}

// Delete performs a logical delete by writing a tombstone.
func (db *DB) Delete(key string) {
	// TODO:
	// 1) panic if db is closed.
	db.mu.Lock()
	if db.mutable == nil {
		panic("DB IS closed")
	}
	// 2) write tombstone into mutable memtable.
	defer db.mu.Unlock()
	db.mutable.Delete(key)
	// 3) rotate with the same rules as Put.
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
