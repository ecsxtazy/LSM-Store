package db

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestStress(t *testing.T) {
	opts := Options{
		MemtableMaxRecords:      100,
		SSTableCompactThreshold: 3,
	}
	db := Open(opts)
	defer db.Close()
	var refMu sync.RWMutex
	refModel := make(map[string][]byte)

	var wg sync.WaitGroup
	stopCh := make(chan struct{})
	numWorkers := 20
	opsPerWorker := 500

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))

			for op := 0; op < opsPerWorker; op++ {
				select {
				case <-stopCh:
					return
				default:
				}

				key := fmt.Sprintf("key_%d_%d", workerID, rng.Intn(100))

				switch rng.Intn(4) {
				case 0, 1:
					value := []byte(fmt.Sprintf("value_%d_%d_%d", workerID, op, rng.Int()))
					db.Put(key, value)

					refMu.Lock()
					copyValue := make([]byte, len(value))
					copy(copyValue, value)
					refModel[key] = copyValue
					refMu.Unlock()

				case 2:
					dbValue, dbOk := db.Get(key)

					refMu.RLock()
					refValue, refOk := refModel[key]
					refMu.RUnlock()

					if dbOk != refOk {
						t.Errorf("Get %s: ok mismatch: db=%v ref=%v", key, dbOk, refOk)
					}
					if dbOk && string(dbValue) != string(refValue) {
						t.Errorf("Get %s: value mismatch: db=%s ref=%s", key, string(dbValue), string(refValue))
					}

				case 3:
					db.Delete(key)

					refMu.Lock()
					delete(refModel, key)
					refMu.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()
	close(stopCh)
}
