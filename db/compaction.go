package db

import "multithreading_hw_3/memtable"

// runCompactor performs background compaction when threshold is exceeded.
func (db *DB) runCompactor() {
	// TODO:
	// 1) wait for compactCh notifications or shutdown signal.
	defer db.wg.Done()
	for {
		select {
		case <-db.closeCh:
			return
		case <-db.compactCh:
			// 2) while len(sstables) > SSTableCompactThreshold, compact oldest pair.
			for {
				db.mu.RLock()
				shouldCompact := len(db.sstables) > db.opts.SSTableCompactThreshold
				db.mu.RUnlock()
				if !shouldCompact {
					break
				}
				if !db.compactOldestPair() {
					break
				}
			}
		}
	}
	// 3) only one compaction at a time.
}

// compactOldestPair merges two oldest SSTables into one and publishes atomically.
// Returns false when compaction is not possible (e.g., fewer than 2 tables).
func (db *DB) compactOldestPair() bool {
	// TODO:
	// 1) snapshot two oldest tables.
	db.mu.Lock()
	if len(db.sstables) < 2 {
		db.mu.Unlock()
		return false
	}
	older := db.sstables[0]
	newest := db.sstables[1]
	// 2) merge outside lock.
	db.mu.Unlock()
	merged := mergeTwoSSTables(older, newest)
	// 3) publish new sstable list atomically (old two removed, merged inserted).
	db.mu.Lock()
	defer db.mu.Unlock()

	newSstables := make([]*SSTable, 0, len(db.sstables)-1)
	found := 0
	for _, sst := range db.sstables {
		if found < 2 && (sst == older || sst == newest) {
			found++
			continue
		}
		newSstables = append(newSstables, sst)
	}
	newSstables = append(newSstables, merged)
	db.sstables = newSstables
	return false
}

// mergeTwoSSTables merges two sorted SSTables where `newer` has higher priority
// for duplicate keys. Tombstones may be retained in the baseline solution.
func mergeTwoSSTables(older, newer *SSTable) *SSTable {
	// TODO: two-way merge by key with version precedence (newer overrides older).
	olderRecords := older.records
	newestRecords := newer.records
	merged := make([]memtable.Record, 0)
	i, j := 0, 0
	for i < len(olderRecords) && j < len(newestRecords) {
		if olderRecords[i].Key < newestRecords[j].Key {
			merged = append(merged, olderRecords[i])
			i++
		} else if olderRecords[i].Key > newestRecords[j].Key {
			merged = append(merged, newestRecords[j])
			j++
		} else {
			merged = append(merged, newestRecords[j])
			i++
			j++
		}
	}
	for i < len(olderRecords) {
		merged = append(merged, olderRecords[i])
		i++
	}
	for j < len(newestRecords) {
		merged = append(merged, newestRecords[j])
		i++
	}
	return &SSTable{
		records: merged,
	}
}
