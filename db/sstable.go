package db

import "multithreading_hw_3/memtable"

// SSTable is an immutable, in-memory sorted list of records.
type SSTable struct {
	// records are sorted by Key ascending and may contain tombstones.
	records []memtable.Record
}

// BuildSSTableFromMemtable converts an immutable memtable into an SSTable.
func BuildSSTableFromMemtable(mt *memtable.Memtable) *SSTable {
	// TODO:
	// 1) iterate immutable memtable in key order.
	// 2) copy each record (including tombstones) into an internal slice.
	// 3) ensure values are copied so SSTable data is immutable from outside.
	panic("TODO")
}

// Get finds a record by key. ok=false means not found.
func (s *SSTable) Get(key string) (memtable.Record, bool) {
	// TODO: binary or linear search by key in sorted records.
	panic("TODO")
}
