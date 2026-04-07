package db

import "multithreading_hw_3/memtable"

func (db *DB) runFlusher() {
	defer db.wg.Done()
	for {
		select {
		case <-db.closeCh:
			for {
				select {
				case immutable := <-db.flushCh:
					sst := BuildSSTableFromMemtable(immutable)
					db.publishFlushed(immutable, sst)
				default:
					return
				}
			}
		case immutable := <-db.flushCh:
			sst := BuildSSTableFromMemtable(immutable)
			db.publishFlushed(immutable, sst)
			db.mu.RLock()
			shouldNotify := len(db.sstables) > db.opts.SSTableCompactThreshold
			db.mu.RUnlock()
			if shouldNotify {
				select {
				case db.compactCh <- struct{}{}:
				default:
				}
			}
		}
	}
}

func (db *DB) publishFlushed(flushed *memtable.Memtable, table *SSTable) {
	db.mu.Lock()
	defer db.mu.Unlock()
	for i, im := range db.immutables {
		if im == flushed {
			db.immutables = append(db.immutables[:i], db.immutables[i+1:]...)
			break
		}
	}
	db.sstables = append(db.sstables, table)
}
