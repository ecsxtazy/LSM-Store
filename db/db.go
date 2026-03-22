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
	panic("TODO")
}

// Close marks DB as closed, flushes pending immutable memtables,
// stops flusher/compactor and waits for their completion.
func (db *DB) Close() {
	// TODO:
	// 1) set closed flag under lock.
	// 2) rotate non-empty mutable into immutables and enqueue for flush.
	// 3) stop background workers and wait for db.wg.
	panic("TODO")
}

// Put writes a value for key into the current mutable memtable.
func (db *DB) Put(key string, value []byte) {
	// TODO:
	// 1) panic if db is closed.
	// 2) write copy(value) into mutable memtable.
	// 3) rotate mutable->immutable after insert when Len reaches MemtableMaxRecords.
	// 4) enqueue immutable into flushCh without blocking forever.
	panic("TODO")
}

// Get searches data in order: mutable -> immutables -> sstables.
func (db *DB) Get(key string) ([]byte, bool) {
	// TODO:
	// 1) panic if db is closed.
	// 2) snapshot pointers under lock.
	// 3) search from newest to oldest across all layers.
	// 4) return copy(value), tombstone means not found.
	panic("TODO")
}

// Delete performs a logical delete by writing a tombstone.
func (db *DB) Delete(key string) {
	// TODO:
	// 1) panic if db is closed.
	// 2) write tombstone into mutable memtable.
	// 3) rotate with the same rules as Put.
	panic("TODO")
}
