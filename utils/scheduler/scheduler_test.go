package scheduler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const MAX_WORKER = 10

func TestScheduler(t *testing.T) {
	dispatcher := NewDispatcher(MAX_WORKER)
	dispatcher.Run()

	for i := 0; i < 100; i++ {
		job := func() {
			t.Log("job is in execution ...")
			time.Sleep(500 * time.Millisecond)
		}

		if ok := dispatcher.AddJob(job); !ok {
			i--
		}
	}
}

func TestSchedulerStop(t *testing.T) {
	dispatcher := NewDispatcher(MAX_WORKER)
	dispatcher.Run()
	dispatcher.Close()
	assert.Equal(t, dispatcher.IsClosed, true)
}
