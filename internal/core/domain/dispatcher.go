package domain

import (
	"github.com/rs/zerolog/log"
)

// Job represents the job to be run
type Job interface {
	Send() error
}

// Worker represents the worker that executes the job
type worker struct {
	workerPool chan chan Job
	jobChannel chan Job
	quit       chan bool
}

func newWorker(workerPool chan chan Job) worker {
	return worker{
		workerPool: workerPool,
		jobChannel: make(chan Job),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w worker) start() {
	go func() {
		log.Info().Msg("starting worker")

		for {
			// register the current worker into the worker queue.
			w.workerPool <- w.jobChannel

			select {
			case job := <-w.jobChannel:
				// we have received a work request.
				if err := job.Send(); err != nil {
					log.Err(err).Msg("Error calling Send")
				}

			case <-w.quit:
				log.Info().Msg("stoping worker")
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w worker) stop() {
	go func() {
		w.quit <- true
	}()
}

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	workerPool chan chan Job

	// A buffered channel that we can send work requests on.
	jobQueue chan Job

	workers []worker
}

func NewDispatcher(maxWorkers, maxQueue int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	jobQueue := make(chan Job, maxQueue)
	return &Dispatcher{workerPool: pool, jobQueue: jobQueue}
}

func (d *Dispatcher) Push(job Job) {
	select {
	case d.jobQueue <- job:
	default:
		log.Warn().Msg("Job Channel full. Discarding value")
	}

}

func (d *Dispatcher) Close() {
	// starting n number of workers
	for _, worker := range d.workers {
		worker.stop()
	}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < cap(d.workerPool); i++ {
		worker := newWorker(d.workerPool)
		worker.start()
		d.workers = append(d.workers, worker)
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for job := range d.jobQueue {
		// a job request has been received
		go func(job Job) {
			// try to obtain a worker job channel that is available.
			// this will block until a worker is idle
			jobChannel := <-d.workerPool

			// dispatch the job to the worker job channel
			jobChannel <- job
		}(job)
	}
}
