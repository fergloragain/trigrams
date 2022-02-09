package main

import (
	"fmt"
	"github.com/fergloragain/trigrams/generate"
	"github.com/fergloragain/trigrams/gram"
	"github.com/fergloragain/trigrams/learn"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// Define consts for now, these should be command line flags
const (
	MaxWorker        = 5
	MaxQueue         = 5
	MaxWords         = 100
	GramSize         = 3
	StripPunctuation = false
)

func main() {

	// If MaxWords is defined, check that it is a reasonable size, i.e. greater than the gram size
	if MaxWords > 0 && MaxWords < GramSize {
		log.Fatal(fmt.Sprintf("Maximum number of words (%d) cannot be less than gram size (%d)", MaxWords, GramSize))
	}

	router := httprouter.New()

	// create the learner queue and workers for handling /learn requests
	learnQueue := make(chan learn.Task, MaxQueue)
	learnDispatcher := learn.NewDispatcher(MaxWorker)
	learnDispatcher.Run(learnQueue, GramSize, StripPunctuation)

	// create the generate queue and workers for handling /generate requests
	generationQueue := make(chan generate.Task, MaxQueue)
	generationDispatcher := generate.NewDispatcher(MaxWorker)
	generationDispatcher.Run(generationQueue, MaxWords, GramSize)

	// the gramCollection is our in-memory data store
	gramCollection := gram.NewCollection()

	// add handlers to the webserver
	handleLearn(router, gramCollection, learnQueue)
	handleGenerate(router, gramCollection, generationQueue)

	// run the server
	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleLearn(router *httprouter.Router, gramCollection *gram.GramCollection, learnQueue chan learn.Task) {
	router.Handle("POST", "/learn", learn.Handler(gramCollection, learnQueue))
}

func handleGenerate(router *httprouter.Router, gramCollection *gram.GramCollection, generationQueue chan generate.Task) {
	router.Handle("GET", "/generate", generate.Handler(gramCollection, generationQueue))
}
