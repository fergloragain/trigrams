package learn

import (
	"fmt"
	"github.com/fergloragain/trigrams/gram"
	"io/ioutil"
	"strings"
	"testing"
)

func TestNewWorker(t *testing.T) {
	workerPool := make(chan chan Task)
	worker := NewWorker(workerPool)

	if worker.LearnWorkerPool != workerPool {
		t.Fail()
	}

	if worker.JobChannel == nil {
		t.Fail()
	}

	if worker.quit == nil {
		t.Fail()
	}
}

func TestStop(t *testing.T) {
	workerPool := make(chan chan Task)
	worker := NewWorker(workerPool)

	jobChannel := make(chan Task)

	worker.JobChannel = jobChannel

	worker.Stop()

	<-worker.quit
}

func TestProcess(t *testing.T) {

	tt := []struct {
		GramSize      int
		Strip         bool
		Gram          *gram.GramCollection
		Regex         []RegexReplacements
		Text          string
		Error         string
		ExpectedGrams [][]string
	}{
		{
			GramSize: 4,
			Strip:    true,
			Gram:     gram.NewCollection(),
			Regex: []RegexReplacements{
				{},
			},
			Text:  "A test input string",
			Error: "",
			ExpectedGrams: [][]string{
				{"A", "test", "input", "string"},
			},
		},
		{
			GramSize: 3,
			Strip:    true,
			Gram:     gram.NewCollection(),
			Regex: []RegexReplacements{
				{},
			},
			Text:  "A test input string",
			Error: "",
			ExpectedGrams: [][]string{
				{"A", "test", "input"},
				{"test", "input", "string"},
			},
		},
		{
			GramSize: 2,
			Strip:    true,
			Gram:     gram.NewCollection(),
			Regex: []RegexReplacements{
				{},
			},
			Text:  "A test input string",
			Error: "",
			ExpectedGrams: [][]string{
				{"A", "test"},
				{"test", "input"},
				{"input", "string"},
			},
		},
		{
			GramSize: 1,
			Strip:    true,
			Gram:     gram.NewCollection(),
			Regex: []RegexReplacements{
				{},
			},
			Text:  "A test input string",
			Error: "",
			ExpectedGrams: [][]string{
				{"A"},
				{"test"},
				{"input"},
				{"string"},
			},
		},
	}

	for _, test := range tt {

		reader := strings.NewReader(test.Text)

		r := ioutil.NopCloser(reader)

		task := &Task{
			Body: r,
			Gram: test.Gram,
			Done: make(chan int),
		}

		var res error

		go func() {
			res = task.Process(test.GramSize, test.Strip, test.Regex)
		}()

		<-task.Done

		if res != nil {
			if res.Error() != test.Error {
				t.Fail()
			}
		}

		for i, g := range test.Gram.Grams {
			for x := range g {
				if test.ExpectedGrams[i][x] != g[x] {
					t.Fail()
				}
			}
		}

	}

}

func TestStripPunctuation(t *testing.T) {

	tt := []struct {
		Text         string
		Replacements []RegexReplacements
		Result       string
		Error        error
	}{
		{
			Text:         "A B",
			Replacements: []RegexReplacements{},
			Result:       "A B",
		},
		{
			Text: "A,B",
			Replacements: []RegexReplacements{
				{
					Regex:       "[^a-zA-Z0-9\\-\\.,!\\?' ]+",
					Replacement: "",
				},
			},
			Result: "A,B",
		},
		{
			Text: "A@B",
			Replacements: []RegexReplacements{
				{
					Regex:       "[^a-zA-Z0-9\\-\\.,!\\?' ]+",
					Replacement: "",
				},
			},
			Result: "AB",
		},
		{
			Text: "A!B",
			Replacements: []RegexReplacements{
				{
					Regex:       "[^a-zA-Z0-9\\-\\.,!\\?' ]+",
					Replacement: "",
				},
			},
			Result: "A!B",
		},
		{
			Text: "Hey how are you?",
			Replacements: []RegexReplacements{
				{
					Regex:       "[^a-zA-Z0-9\\-\\.,!\\?' ]+",
					Replacement: "",
				},
			},
			Result: "Hey how are you?",
		},
		{
			Text: `This is
a test`,
			Replacements: []RegexReplacements{
				{
					Regex:       `\r?\n`,
					Replacement: " ",
				},
			},
			Result: "This is a test",
		},
		{
			Text: `Someone
call
999!`,
			Replacements: []RegexReplacements{
				{
					Regex:       `\r?\n`,
					Replacement: " ",
				},
			},
			Result: "Someone call 999!",
		},
		{
			Text: `Let's eat, Grandma.`,
			Replacements: []RegexReplacements{
				{
					Regex:       `\r?\n`,
					Replacement: " ",
				},
				{
					Regex:       "[^a-zA-Z0-9\\-\\.,!\\?' ]+",
					Replacement: "",
				},
			},
			Result: "Let's eat, Grandma.",
		},
		{
			Text: `Let's eat Grandma!`,
			Replacements: []RegexReplacements{
				{
					Regex:       `\r?\n`,
					Replacement: " ",
				},
				{
					Regex:       "[^a-zA-Z0-9\\-\\.,!\\?' ]+",
					Replacement: "",
				},
			},
			Result: "Let's eat Grandma!",
		},
		{
			Text: `me@gmail.com`,
			Replacements: []RegexReplacements{
				{
					Regex:       "[^a-zA-Z0-9\\-\\.,!\\?' ]+",
					Replacement: "",
				},
			},
			Result: "megmail.com",
		},
		{
			Text: `this will fail to compile`,
			Replacements: []RegexReplacements{
				{
					Regex:       "[.-()]",
					Replacement: "",
				},
			},
			Result: "",
		},
	}

	for _, test := range tt {
		res, _ := stripPunctuation(test.Text, test.Replacements)

		if res != test.Result {
			t.Errorf(fmt.Sprintf("Expected %s to be %s", res, test.Result))
		}

	}

}

func TestProcessToTokens(t *testing.T) {

	tt := []struct {
		Text      string
		GramSize  int
		Result    []string
		Remaining string
		Error     error
	}{
		{
			Text:      "A B",
			GramSize:  1,
			Result:    []string{"A"},
			Remaining: "B",
		},
		{
			Text:      "A B C",
			GramSize:  2,
			Result:    []string{"A", "B"},
			Remaining: "B C",
		},
		{
			Text:      "A B C ",
			GramSize:  2,
			Result:    []string{"A", "B"},
			Remaining: "B C ",
		},
		{
			Text:      "A B C D",
			GramSize:  3,
			Result:    []string{"A", "B", "C"},
			Remaining: "B C D",
		},
		{
			Text:      "A B C D E ",
			GramSize:  4,
			Result:    []string{"A", "B", "C", "D"},
			Remaining: "B C D E ",
		},
		{
			Text:      "something more than individual letters ",
			GramSize:  4,
			Result:    []string{"something", "more", "than", "individual"},
			Remaining: "more than individual letters ",
		},
		{
			Text: `
a 
little
more
complex piece of text
with line


breaks`,
			GramSize:  4,
			Result:    []string{"a", "little", "more", "complex"},
			Remaining: "little more complex piece of text with line breaks",
		},
	}

	for _, test := range tt {
		res, remaining := processToTokens(test.Text, test.GramSize)

		if len(res) != len(test.Result) {
			fmt.Println(res)
			fmt.Println(test.Result)
		}
		for x, v := range test.Result {
			if v != res[x] {
				t.Fail()
			}
		}

		if remaining != test.Remaining {
			t.Errorf(fmt.Sprintf("Expected >%s< to be >%s<", remaining, test.Remaining))
		}

	}

}
