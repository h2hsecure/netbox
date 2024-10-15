package domain

import (
	"context"

	"github.com/rs/zerolog/log"
)

type Job interface {
	Send(context.Context) error
}

type worker struct {
	workerPool chan chan Job
	jobChannel chan Job
	quit       chan bool
	ctx        context.Context
}

func newWorker(workerPool chan chan Job) *worker {
	return &worker{
		workerPool: workerPool,
		jobChannel: make(chan Job),
		quit:       make(chan bool),
		ctx:        context.Background()}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w *worker) start() {
	go func() {
		log.Info().Msg("starting worker")

		for {
			// register the current worker into the worker queue.
			w.workerPool <- w.jobChannel

			select {
			case job := <-w.jobChannel:
				// we have received a work request.
				if err := job.Send(w.ctx); err != nil {
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
func (w *worker) stop() {
	go func() {
		w.quit <- true
	}()
}

type Dispatcher struct {
	workerPool chan chan Job

	jobQueue chan Job

	workers []*worker
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
	for _, worker := range d.workers {
		worker.stop()
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < cap(d.workerPool); i++ {
		worker := newWorker(d.workerPool)
		worker.start()
		d.workers = append(d.workers, worker)
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for job := range d.jobQueue {
		go func(job Job) {
			jobChannel := <-d.workerPool
			jobChannel <- job
		}(job)
	}
}
