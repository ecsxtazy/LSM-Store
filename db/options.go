package db

// Options controls DB behavior.
type Options struct {
	// MemtableMaxRecords triggers rotation when Len reaches this value.
	MemtableMaxRecords int
	// SSTableCompactThreshold triggers compaction when len(sstables) is greater than this value.
	SSTableCompactThreshold int
}

// DefaultOptions returns a minimal set of sane defaults.
func DefaultOptions() Options {
	return Options{
		MemtableMaxRecords:      1024,
		SSTableCompactThreshold: 4,
	}
}
