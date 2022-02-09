package generate

import (
	"github.com/fergloragain/trigrams/gram"
	"log"
	"net/http"
)

type Task struct {
	Writer http.ResponseWriter
	Gram   *gram.GramCollection
	Output chan string
}

type GenerationWorker struct {
	WorkerPool        chan chan Task
	GenerationChannel chan Task
	quit              chan bool
}

func NewGenerationWorker(workerPool chan chan Task) GenerationWorker {
	return GenerationWorker{
		WorkerPool:        workerPool,
		GenerationChannel: make(chan Task),
		quit:              make(chan bool),
	}
}

// start the generation worker, specifying the maximum number of words to be generated, and the size of the ngram we are
// working with
func (w GenerationWorker) Start(maxWords, gramSize int) {
	go func() {
		for {
			// the worker registers itself into the pool of workers
			w.WorkerPool <- w.GenerationChannel

			select {
			// the worker listens for a generationTask request
			case generationTask := <-w.GenerationChannel:

				// process the generationTask request
				randomText, err := generationTask.Process(maxWords, gramSize)

				if err != nil {
					log.Printf("Error generating text: %s", err.Error())
				}

				// write the random text to the output channel
				generationTask.Output <- randomText

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w GenerationWorker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func (task *Task) Process(max, gramSize int) (string, error) {
	// build random text based on the grams that have been learned
	randomString, err := task.Gram.BuildRandomText(max, gramSize)

	if err != nil {
		return "", err
	}

	return randomString, nil
}
