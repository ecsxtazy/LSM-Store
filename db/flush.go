package db

import "multithreading_hw_3/memtable"

// runFlusher consumes immutable memtables and builds SSTables in FIFO order.
func (db *DB) runFlusher() {
	// TODO:
	// 1) receive immutable memtables from flushCh until shutdown.
	// 2) build SSTable from each memtable.
	// 3) publish result atomically: add new SSTable and remove immutable from visible list.
	// 4) notify compactor when sstable count may exceed threshold.
	panic("TODO")
}

// publishFlushed atomically replaces state visible to Get/Range after flush.
func (db *DB) publishFlushed(flushed *memtable.Memtable, table *SSTable) {
	// TODO: update db.immutables/db.sstables in one short critical section.
	panic("TODO")
}
