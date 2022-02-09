package learn

import (
	"github.com/fergloragain/trigrams/gram"
	"io/ioutil"
	"strings"
	"testing"
)

func TestNewDispatcher(t *testing.T) {
	d := NewDispatcher(1)

	if d.WorkerPool == nil {
		t.Fail()
	}

	if d.numberOfWorkers != 1 {
		t.Fail()
	}
}

func TestRun(t *testing.T) {
	d := NewDispatcher(1)

	if d.WorkerPool == nil {
		t.Fail()
	}

	if d.numberOfWorkers != 1 {
		t.Fail()
	}

	learnQueue := make(chan Task)

	go d.Run(learnQueue, 1, false)

	reader := strings.NewReader("")

	r := ioutil.NopCloser(reader)

	learnQueue <- Task{
		Body: r,
		Gram: gram.NewCollection(),
	}

}
