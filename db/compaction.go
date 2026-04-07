package db

import "multithreading_hw_3/memtable"

func (db *DB) runCompactor() {
	defer db.wg.Done()
	for {
		select {
		case <-db.closeCh:
			return
		case <-db.compactCh:
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
}
func (db *DB) compactOldestPair() bool {
	db.mu.Lock()
	if len(db.sstables) < 2 {
		db.mu.Unlock()
		return false
	}

	older := db.sstables[0]
	newer := db.sstables[1]

	olderRecords := make([]memtable.Record, len(older.records))
	copy(olderRecords, older.records)
	newerRecords := make([]memtable.Record, len(newer.records))
	copy(newerRecords, newer.records)

	db.mu.Unlock()
	merged := mergeTwoSSTables(&SSTable{records: olderRecords}, &SSTable{records: newerRecords})
	db.mu.Lock()
	defer db.mu.Unlock()

	if len(db.sstables) < 2 || db.sstables[0] != older || db.sstables[1] != newer {
		return false
	}

	newSstables := make([]*SSTable, 0, len(db.sstables)-1)
	newSstables = append(newSstables, merged)
	newSstables = append(newSstables, db.sstables[2:]...)
	db.sstables = newSstables

	return true
}

func mergeTwoSSTables(older, newer *SSTable) *SSTable {
	olderRecords := older.records
	newerRecords := newer.records
	merged := make([]memtable.Record, 0)
	i, j := 0, 0
	for i < len(olderRecords) && j < len(newerRecords) {
		if olderRecords[i].Key < newerRecords[j].Key {
			merged = append(merged, copyRecord(olderRecords[i]))
			i++
		} else if olderRecords[i].Key > newerRecords[j].Key {
			merged = append(merged, copyRecord(newerRecords[j]))
			j++
		} else {
			merged = append(merged, copyRecord(newerRecords[j]))
			i++
			j++
		}
	}
	for i < len(olderRecords) {
		merged = append(merged, copyRecord(olderRecords[i]))
		i++
	}
	for j < len(newerRecords) {
		merged = append(merged, copyRecord(newerRecords[j]))
		j++
	}
	return &SSTable{
		records: merged,
	}
}

func copyRecord(rec memtable.Record) memtable.Record {
	copied := memtable.Record{
		Key:       rec.Key,
		Tombstone: rec.Tombstone,
	}
	if rec.Value != nil {
		copied.Value = make([]byte, len(rec.Value))
		copy(copied.Value, rec.Value)
	}
	return copied
}
