package db

// Range iterates visible keys in [start, end) in ascending key order.
// start == "" means no left bound, end == "" means no right bound.
// If start >= end (and both are non-empty), Range returns immediately.
func (db *DB) Range(start, end string, fn func(key string, value []byte) bool) {
	// TODO:
	// 1) panic if db is closed.
	// 2) take a consistent snapshot of mutable/immutables/sstables under lock.
	// 3) merge newest->oldest versions, keep only first version per key.
	// 4) drop tombstones, sort keys ascending, invoke fn with copy(value).
	// 5) stop immediately if fn returns false.
	panic("TODO")
}
