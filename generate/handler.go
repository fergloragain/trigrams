package generate

import (
	"fmt"
	"github.com/fergloragain/trigrams/gram"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Handler(gram *gram.GramCollection, generationQueue chan Task) func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		output := make(chan string)

		generationJob := Task{
			Writer: writer,
			Gram:   gram,
			Output: output,
		}

		generationQueue <- generationJob

		generatedText := <-output

		fmt.Fprintf(writer, generatedText)
	}

}
