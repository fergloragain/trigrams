package learn

import (
	"github.com/fergloragain/trigrams/gram"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Handler(gram *gram.GramCollection, learnQueue chan Task) func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	return func(writer http.ResponseWriter, request *http.Request, p httprouter.Params) {
		job := Task{
			Body: request.Body,
			Gram: gram,
			Done: make(chan int),
		}

		learnQueue <- job

		<-job.Done

		writer.WriteHeader(http.StatusOK)
	}
}
