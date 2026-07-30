package memtable

import (
	"slices"
	"sort"
	"sync"
)

type Record struct {
	Key       string
	Value     []byte
	Tombstone bool
}

type Memtable struct {
	mu   sync.RWMutex
	data []Record
}

func New() *Memtable {
	return &Memtable{
		data: make([]Record, 0),
	}
}
func (m *Memtable) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

func (m *Memtable) Get(key string) ([]byte, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	index := sort.Search(len(m.data), func(i int) bool {
		return m.data[i].Key >= key
	})
	if index < len(m.data) && m.data[index].Key == key {
		if m.data[index].Tombstone {
			return nil, false
		}
		return append([]byte(nil), m.data[index].Value...), true
	}
	return nil, false
}

func (m *Memtable) Put(key string, value []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	index := sort.Search(len(m.data), func(i int) bool {
		return m.data[i].Key >= key
	})
	copy := append([]byte(nil), value...)
	if index < len(m.data) && m.data[index].Key == key {
		m.data[index].Value = copy
		m.data[index].Tombstone = false
	} else {
		m.data = slices.Insert(m.data, index, Record{
			Key:       key,
			Value:     copy,
			Tombstone: false,
		})
	}
}

func (m *Memtable) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	index := sort.Search(len(m.data), func(i int) bool {
		return m.data[i].Key >= key
	})
	if index < len(m.data) && m.data[index].Key == key {
		m.data[index].Tombstone = true
	} else {
		m.data = slices.Insert(m.data, index, Record{
			Key:       key,
			Value:     nil,
			Tombstone: true,
		})
	}
}

func (m *Memtable) Range(fn func(rec Record) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, rec := range m.data {
		if !fn(rec) {
			break
		}
	}
}
