package learn

import (
	"fmt"
	"github.com/fergloragain/trigrams/gram"
	"io"
	"log"
	"regexp"
	"strings"
)

const ReadSize = 64

var regexReplacements []RegexReplacements

type Task struct {
	Body io.ReadCloser
	Gram *gram.GramCollection
	Done chan int
}

type LearnWorker struct {
	LearnWorkerPool chan chan Task
	JobChannel      chan Task
	quit            chan bool
}

func NewWorker(workerPool chan chan Task) LearnWorker {
	return LearnWorker{
		LearnWorkerPool: workerPool,
		JobChannel:      make(chan Task),
		quit:            make(chan bool),
	}
}

type RegexReplacements struct {
	Regex       string
	Replacement string
}

func init() {
	regexReplacements = []RegexReplacements{}

	regexReplacements = append(regexReplacements, RegexReplacements{
		Regex:       `\r?\n`,
		Replacement: " ",
	})

	regexReplacements = append(regexReplacements, RegexReplacements{
		Regex:       "[^a-zA-Z0-9\\-\\.,!\\?' ]+",
		Replacement: "",
	})
}

// start the learn worker, specifying the maximum number of words to be generated, and whether or not we want to strip
// punctuation
func (worker LearnWorker) Start(gramSize int, stripPunctuation bool) {
	go func() {
		for {
			// the worker registers itself into the pool of workers
			worker.LearnWorkerPool <- worker.JobChannel

			select {
			// the worker listens for a learnTask request
			case learnTask := <-worker.JobChannel:

				// process the learnTask request
				if err := learnTask.Process(gramSize, stripPunctuation, regexReplacements); err != nil {
					log.Printf("Error processing job: %s", err.Error())
				}

			case <-worker.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (worker LearnWorker) Stop() {
	go func() {
		worker.quit <- true
	}()
}

// Process will strip punctuation from the source text if configured to do so, then split the text into an array of
// strings, and then process the array of strings into ngrams of a specific size, by default 3
func (job *Task) Process(gramSize int, strip bool, regexArray []RegexReplacements) error {

	defer job.Body.Close()

	streamBuffer := make([]byte, ReadSize)

	var remainingString string
	gramTokens := []string{}

	for {

		numberOfBytesRead, err := job.Body.Read(streamBuffer)

		if numberOfBytesRead > 0 {

			byteString := string(streamBuffer[:numberOfBytesRead])

			var plainText string

			plainText, err = stripPunctuation(byteString, regexArray)

			if err != nil {
				return err
			}

			if len(plainText) > 0 {

				plainText = fmt.Sprintf("%s%s", remainingString, plainText)

				gramTokens, remainingString = processToTokens(plainText, gramSize)

				for len(gramTokens) == gramSize {

					newGram := []string{}

					newGram = append(newGram, gramTokens[:gramSize]...)

					job.Gram.AddGram(newGram)

					gramTokens, remainingString = processToTokens(remainingString, gramSize)

				}
			}
		} else {
			break
		}
	}

	if len(remainingString) > 0 {
		lastGram := strings.Fields(remainingString)

		if len(lastGram) == gramSize {
			job.Gram.AddGram(lastGram)
		}
	}

	job.Done <- 1

	return nil
}

// processToTokens accepts a string and a gram size, and attempts to parse the string into a gram of the corresponding
// size. If there isn't enough words in the string to create a gram of the correct size, an empty gram is returned
// along with the input string, which will be processed with more data later on
func processToTokens(fullString string, gramSize int) ([]string, string) {

	tokens := strings.Fields(fullString)

	if len(tokens) > gramSize {

		gramTokens := tokens[:gramSize]

		trailingSize := gramSize - 1

		remaining := []string{}

		remaining = append(remaining, gramTokens[len(gramTokens)-trailingSize:]...)

		remaining = append(remaining, tokens[gramSize:]...)

		remainingString := strings.Join(remaining, " ")

		if strings.HasSuffix(fullString, " ") {
			remainingString = fmt.Sprintf("%s ", remainingString)
		}

		return gramTokens, remainingString
	}

	return []string{}, fullString
}

// stripPunctuation accepts a string and an array of regex patterns and replacement strings, and modifies the input
// string by replacing each regex pattern with the corresponding pattern
func stripPunctuation(text string, regexReplacements []RegexReplacements) (string, error) {

	for _, regexReplacement := range regexReplacements {
		regex, err := regexp.Compile(regexReplacement.Regex)

		if err != nil {
			return "", err
		}

		text = regex.ReplaceAllString(text, regexReplacement.Replacement)
	}

	return text, nil
}
