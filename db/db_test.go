package db

import (
	"testing"
)

func TestRangeReturnsSortedVisibleKeys(t *testing.T) {
	t.Skip("TODO: implement test")

	// Outline:
	// 1) Put keys so data is split across mutable/immutable/sstables.
	// 2) Call Range("", "", fn) and collect keys.
	// 3) Verify strict ascending order.
	// 4) Verify no duplicates.
}

func TestRangeSkipsTombstonesAndKeepsNewestVersion(t *testing.T) {
	t.Skip("TODO: implement test")

	// Outline:
	// 1) Flush old value into SSTable.
	// 2) Write newer tombstone or newer value in mutable/immutable.
	// 3) Range over key interval.
	// 4) Verify tombstoned keys are absent and newest value wins.
}

func TestCompactionReplacesTwoOldestSSTables(t *testing.T) {
	t.Skip("TODO: implement test")

	// Outline:
	// 1) Create multiple SSTables with overlapping keys and small compact threshold.
	// 2) Wait for compaction of two oldest tables.
	// 3) Verify resulting active SSTables count/ordering is valid.
	// 4) Verify Get/Range values are preserved.
}

func TestRangeStopsWhenCallbackReturnsFalse(t *testing.T) {
	t.Skip("TODO: implement test")

	// Outline:
	// 1) Fill DB with several keys.
	// 2) Range with callback that returns false after N items.
	// 3) Verify callback is not called for later keys.
}

func TestCloseStopsFlusherAndCompactor(t *testing.T) {
	t.Skip("TODO: implement test")

	// Outline:
	// 1) Start DB, force flush and compaction activity.
	// 2) Call Close().
	// 3) Verify no background goroutines keep running (via WaitGroup-visible state).
	// 4) Verify operations started after close panic.
}
