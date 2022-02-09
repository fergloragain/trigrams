package learn

import (
	"bytes"
	"github.com/fergloragain/trigrams/gram"
	"github.com/pkg/errors"
	"net/http"
	"testing"
)

type TestWriter struct {
	ResultCode int
	Blocker    chan int
}

func (t *TestWriter) Header() http.Header {
	return nil
}

func (t *TestWriter) Write(b []byte) (int, error) {
	return 0, nil
}

func (t *TestWriter) WriteHeader(statusCode int) {
	t.ResultCode = statusCode
	t.Blocker <- 0

	return
}

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() error {
	return nil
}

func TestHandler(t *testing.T) {
	gramCollection := gram.NewCollection()
	learnQueue := make(chan Task)

	handler := Handler(gramCollection, learnQueue)

	testWriter := &TestWriter{}
	testWriter.Blocker = make(chan int)

	cb := &ClosingBuffer{bytes.NewBufferString("Hello")}

	testRequest := &http.Request{
		Body: cb,
	}

	go handler(testWriter, testRequest, nil)

	r := <-learnQueue

	r.Done <- 1

	<-testWriter.Blocker

	buf := make([]byte, 6)

	n, _ := r.Body.Read(buf)

	if string(buf[:n]) != "Hello" {
		t.Fail()
	}

	if testWriter.ResultCode != http.StatusOK {
		t.Fail()
	}
}

type BrokenBuffer struct {
	*bytes.Buffer
}

func (cb *BrokenBuffer) Close() error {
	return errors.New("Error reading data")
}

func (cb *BrokenBuffer) Read(p []byte) (n int, err error) {
	return 0, errors.New("Error reading data")
}

type TestWriterError struct {
	Result  string
	Blocker chan int
}

func (t *TestWriterError) Header() http.Header {
	return nil
}

func (t *TestWriterError) Write(b []byte) (int, error) {
	t.Result = string(b)
	t.Blocker <- 0

	return 0, nil
}

func (t *TestWriterError) WriteHeader(statusCode int) {
	return
}
