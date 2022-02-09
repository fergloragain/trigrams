package learn

type LearnDispatcher struct {
	numberOfWorkers int
	WorkerPool      chan chan Task
}

// Create a generation dispatcher, specifying a number of workers to read from the queue of learn tasks
func NewDispatcher(maxWorkers int) *LearnDispatcher {
	pool := make(chan chan Task, maxWorkers)
	return &LearnDispatcher{WorkerPool: pool, numberOfWorkers: maxWorkers}
}

// Run the learn dispatcher by launching each worker and then listening to the learn queue
func (dispatcher *LearnDispatcher) Run(learnQueue chan Task, gramSize int, strip bool) {
	// starting n number of workers
	for i := 0; i < dispatcher.numberOfWorkers; i++ {
		worker := NewWorker(dispatcher.WorkerPool)
		worker.Start(gramSize, strip)
	}

	go dispatcher.dispatch(learnQueue)
}

// listen to the learn queue and once a learn request comes in, retrieve a worker from the worker pool and
// hand the request off to the worker for processing
func (dispatcher *LearnDispatcher) dispatch(learnQueue chan Task) {
	for {
		select {
		// listen for a learn request
		case learnTask := <-learnQueue:
			go func(task Task) {
				// obtain a worker from the worker pool
				learnWorker := <-dispatcher.WorkerPool

				// dispatch the job to the worker job channel
				learnWorker <- task
			}(learnTask)
		}
	}
}
