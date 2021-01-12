package scheduler

import (
	"runtime"
	"time"

	"github.com/meshplus/gosdk/common"
)

var (
	logger = common.GetLogger("scheduler")
)

// Job scheduler must implement method
type Job func()

// Dispatcher is used to dispatch job
type Dispatcher struct {
	TaskChannel     chan Job      // Store all job
	DispatchControl chan bool     // Control add job go routine
	WorkerPool      chan chan Job // Store worker private job channel
	WorkerList      []*Worker     // Store worker instance
	MaxWorker       int           // Max worker number
	IsClosed        bool          // Dispatcher is closed ?
}

// Worker struct
type Worker struct {
	Dispatcher *Dispatcher
	JobChannel chan Job
	Quit       chan bool
	ID         int
}

// NewWorker is used to create a new worker
func NewWorker(d *Dispatcher, ID int) *Worker {
	return &Worker{
		Dispatcher: d,
		JobChannel: make(chan Job),
		Quit:       make(chan bool),
		ID:         ID,
	}
}

// Start worker
func (w *Worker) Start() {
	logger.Info("worker start ", w.ID)
	go func() {
		for {
			w.Dispatcher.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				logger.Debug("Scheduler worker start deal job ...")
				job()
				<-w.Dispatcher.DispatchControl
			case <-w.Quit:
				logger.Debugf("Scheduler worker %d quit.", w.ID)
				return
			}
		}
	}()
}

// Stop worker
func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

// NewDispatcher is used to create dispatch
func NewDispatcher(maxWorker int) *Dispatcher {
	return &Dispatcher{
		TaskChannel:     make(chan Job),
		WorkerPool:      make(chan chan Job, maxWorker),
		DispatchControl: make(chan bool, 2*maxWorker),
		MaxWorker:       maxWorker,
	}
}

func (d *Dispatcher) dispatch() {
	for job := range d.TaskChannel {
		go func(job Job) {
			jobChannel := <-d.WorkerPool
			jobChannel <- job
		}(job)
	}
}

// Run dispatcher
func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorker; i++ {
		worker := NewWorker(d, i)
		d.WorkerList = append(d.WorkerList, worker)
		worker.Start()
	}
	d.IsClosed = false
	go d.dispatch()
}

// Close dispatch
func (d *Dispatcher) Close() {
	d.IsClosed = true
	go func() {
		for _, worker := range d.WorkerList {
			worker.Stop()
		}
	}()
}

// AddJob is used job to global job channel
func (d *Dispatcher) AddJob(job Job) bool {
	select {
	case d.DispatchControl <- true:
		if d.IsClosed {
			logger.Fatalf("Dispatcher is closed, can not add job into dispatcher.")
		}
		d.TaskChannel <- job
		return true
	case <-time.After(time.Microsecond * 100):
		logger.Debug("GoSDK scheduler pool too busy.")
		logger.Debug("goroutine num ->", runtime.NumGoroutine())
		return false
	}
}
