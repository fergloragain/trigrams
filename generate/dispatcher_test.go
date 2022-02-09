package generate

import (
	"fmt"
	"github.com/fergloragain/trigrams/gram"
	"testing"
)

func TestNewDispatcher(t *testing.T) {
	d := NewDispatcher(1)

	if d.WorkerPool == nil {
		t.Fail()
	}

	if d.maxWorkers != 1 {
		t.Fail()
	}
}

func TestRun(t *testing.T) {
	d := NewDispatcher(1)

	if d.WorkerPool == nil {
		t.Fail()
	}

	if d.maxWorkers != 1 {
		t.Fail()
	}

	learnQueue := make(chan Task)

	go d.Run(learnQueue, 1, 3)

	o := make(chan string)

	learnQueue <- Task{
		Writer: nil,
		Gram:   gram.NewCollection(),
		Output: o,
	}

	x := <-o

	fmt.Println(x)

}
