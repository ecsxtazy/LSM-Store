package memtable

type Record struct {
	Key       string
	Value     []byte
	Tombstone bool
}

type Memtable struct {
	// TODO: define internal fields (storage + sync)
}

func New() *Memtable {
	panic("not implemented")
}

func (m *Memtable) Len() int {
	panic("not implemented")
}

func (m *Memtable) Get(key string) ([]byte, bool) {
	panic("not implemented")
}

func (m *Memtable) Put(key string, value []byte) {
	panic("not implemented")
}

func (m *Memtable) Delete(key string) {
	panic("not implemented")
}

func (m *Memtable) Range(fn func(rec Record) bool) {
	panic("not implemented")
}
