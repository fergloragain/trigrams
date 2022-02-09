package generate

type GenerationDispatcher struct {
	maxWorkers int
	WorkerPool chan chan Task
}

// Create a generation dispatcher, specifying a number of workers to read from the queue of generation tasks
func NewDispatcher(maxWorkers int) *GenerationDispatcher {
	pool := make(chan chan Task, maxWorkers)
	return &GenerationDispatcher{WorkerPool: pool, maxWorkers: maxWorkers}
}

// Run the generation dispatcher by launching each worker and then listening to the generation queue
func (dispatcher *GenerationDispatcher) Run(generationQueue chan Task, max, gramSize int) {
	for i := 0; i < dispatcher.maxWorkers; i++ {
		worker := NewGenerationWorker(dispatcher.WorkerPool)
		worker.Start(max, gramSize)
	}

	go dispatcher.dispatch(generationQueue)
}

// listen to the generation queue and once a generation request comes in, retrieve a worker from the worker pool and
// hand the request off to the worker for processing
func (dispatcher *GenerationDispatcher) dispatch(generationQueue chan Task) {
	for {

		select {
		// listen for a generation request
		case generationRequest := <-generationQueue:
			go func(generationRequest Task) {
				// obtain a worker from the worker pool
				generationChannel := <-dispatcher.WorkerPool

				// dispatch the job to the worker job channel
				generationChannel <- generationRequest
			}(generationRequest)
		}
	}
}
