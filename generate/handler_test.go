package generate

import (
	"github.com/fergloragain/trigrams/gram"
	"net/http"
	"testing"
)

type TestWriter struct {
	Result  string
	Blocker chan int
}

func (t *TestWriter) Header() http.Header {
	return nil
}

func (t *TestWriter) Write(b []byte) (int, error) {
	t.Result = string(b)
	t.Blocker <- 0

	return 0, nil
}

func (t *TestWriter) WriteHeader(statusCode int) {
	return
}

func TestHandler(t *testing.T) {
	gramCollection := gram.NewCollection()
	generationQueue := make(chan Task)

	handler := Handler(gramCollection, generationQueue)

	testWriter := &TestWriter{}
	testWriter.Blocker = make(chan int)

	go handler(testWriter, nil, nil)

	r := <-generationQueue

	go func() {
		r.Output <- "TEST"
	}()

	<-testWriter.Blocker

	if testWriter.Result != "TEST" {
		t.Fail()
	}
}
