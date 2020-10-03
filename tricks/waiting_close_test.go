package tricks

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestWorker_WaitingClose(t *testing.T) {

	var (
		executedCount int
		count         = 1000
	)
	w := NewWorker(100, func(job int) {
		time.Sleep(time.Microsecond * 10) // emulate work
		executedCount++
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < count; i++ {
			w.AddJob(i)
		}
	}()
	wg.Wait()

	w.Close()

	assert.Equal(t, count, executedCount)
}
