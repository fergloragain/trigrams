package generate

import (
	"fmt"
	"github.com/fergloragain/trigrams/gram"
	"github.com/pkg/errors"
	"testing"
)

func TestStop(t *testing.T) {
	workerPool := make(chan chan Task)
	worker := NewGenerationWorker(workerPool)

	generationChannel := make(chan Task)

	worker.GenerationChannel = generationChannel

	worker.Stop()

	<-worker.quit
}

func TestFulfil(t *testing.T) {

	gc1 := gram.NewCollection()
	gc1.Grams = [][]string{
		[]string{
			"this", "is", "cool",
		},
	}
	gc1.Indices = []int{0}
	gc1.Frequencies = []int{1}
	gc1.TotalFrequencies = 1

	tt := []struct {
		Task     Task
		Result   string
		Error    error
		Max      int
		GramSize int
	}{
		{
			Task: Task{
				Writer: nil,
				Gram:   gram.NewCollection(),
				Output: nil,
			},
			Result:   "",
			Error:    errors.New("No grams to fetch randomly"),
			Max:      1,
			GramSize: 3,
		},
		{
			Task: Task{
				Writer: nil,
				Gram:   gc1,
				Output: nil,
			},
			Result:   "this is cool",
			Error:    nil,
			Max:      100,
			GramSize: 3,
		},
	}

	for _, tc := range tt {
		str, err := tc.Task.Process(tc.Max, tc.GramSize)

		if str != tc.Result {
			t.Error(fmt.Sprintf("Expected >%s< to match >%s<", tc.Result, str))
		}

		if tc.Error != nil {
			if err.Error() != tc.Error.Error() {
				t.Fail()
			}
		}
	}

}
