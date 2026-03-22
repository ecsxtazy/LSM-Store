package db

// runCompactor performs background compaction when threshold is exceeded.
func (db *DB) runCompactor() {
	// TODO:
	// 1) wait for compactCh notifications or shutdown signal.
	// 2) while len(sstables) > SSTableCompactThreshold, compact oldest pair.
	// 3) only one compaction at a time.
	panic("TODO")
}

// compactOldestPair merges two oldest SSTables into one and publishes atomically.
// Returns false when compaction is not possible (e.g., fewer than 2 tables).
func (db *DB) compactOldestPair() bool {
	// TODO:
	// 1) snapshot two oldest tables.
	// 2) merge outside lock.
	// 3) publish new sstable list atomically (old two removed, merged inserted).
	panic("TODO")
}

// mergeTwoSSTables merges two sorted SSTables where `newer` has higher priority
// for duplicate keys. Tombstones may be retained in the baseline solution.
func mergeTwoSSTables(older, newer *SSTable) *SSTable {
	// TODO: two-way merge by key with version precedence (newer overrides older).
	panic("TODO")
}
