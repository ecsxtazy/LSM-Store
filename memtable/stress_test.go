package memtable

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestStress(t *testing.T) {
	if testing.Short() {
		t.Skip("stress test skipped with -short")
	}

	mt := New()
	const workers = 8
	const ops = 20000

	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func(id int) {
			defer wg.Done()
			rnd := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))
			for i := 0; i < ops; i++ {
				key := fmt.Sprintf("k%04d", rnd.Intn(1000))
				switch rnd.Intn(3) {
				case 0:
					mt.Put(key, []byte("v"))
				case 1:
					mt.Get(key)
				default:
					mt.Delete(key)
				}
			}
		}(w)
	}
	wg.Wait()
}
