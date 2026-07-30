package db

type Options struct {
	MemtableMaxRecords      int
	SSTableCompactThreshold int
}

func DefaultOptions() Options {
	return Options{
		MemtableMaxRecords:      1024,
		SSTableCompactThreshold: 4,
	}
}
