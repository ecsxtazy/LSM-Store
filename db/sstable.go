package db

import (
	"multithreading_hw_3/memtable"
	"sort"
)

type SSTable struct {
	records []memtable.Record
}

func BuildSSTableFromMemtable(mt *memtable.Memtable) *SSTable {
	records := make([]memtable.Record, 0, mt.Len())
	mt.Range(func(rec memtable.Record) bool {
		recCopy := memtable.Record{
			Key:       rec.Key,
			Tombstone: rec.Tombstone,
		}
		if rec.Value != nil {
			recCopy.Value = make([]byte, len(rec.Value))
			copy(recCopy.Value, rec.Value)
		}
		records = append(records, recCopy)
		return true
	})
	return &SSTable{records: records}
}

func (s *SSTable) Get(key string) (memtable.Record, bool) {
	idx := sort.Search(len(s.records), func(i int) bool {
		return s.records[i].Key >= key
	})
	if idx < len(s.records) && s.records[idx].Key == key {
		rec := s.records[idx]
		if rec.Tombstone {
			return memtable.Record{Tombstone: true}, true
		}
		if rec.Value != nil {
			valueCopy := make([]byte, len(rec.Value))
			copy(valueCopy, rec.Value)
			return memtable.Record{Key: rec.Key, Value: valueCopy}, true
		}
		return rec, true
	}
	return memtable.Record{}, false
}
