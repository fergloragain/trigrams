package gram

import (
	"github.com/pkg/errors"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type GramCollection struct {
	Grams            [][]string
	Frequencies      []int
	Indices          []int
	RW               sync.RWMutex
	TotalFrequencies int
}

// Creates a new collection
func NewCollection() *GramCollection {
	grams := new(GramCollection)
	grams.Grams = [][]string{}  // Grams is an array of string arrays
	grams.Frequencies = []int{} // Frequencies is an array of ints corresponding to each Gram's frequency across all learned texts
	grams.Indices = []int{}     // indices is an array of ints specifying the order in which the Grams should be iterated
	return grams
}

// Adding a new gram involves append a gram ([]string) to a 2D array, as well as adding a frequency of 1 (since this will
// be the first time we've encountered this specific gram), as well as a default index, and then incrementing the total
// frequencies of all grams across all learned texts
func (gramCollection *GramCollection) addNewGram(newNgram []string) {
	gramCollection.Grams = append(gramCollection.Grams, newNgram)
	gramCollection.Frequencies = append(gramCollection.Frequencies, 1)
	gramCollection.Indices = append(gramCollection.Indices, len(gramCollection.Indices))
	gramCollection.TotalFrequencies += 1
}

// Updating a frequency increments a frequency at a given index, gramIndex, and then increments the total frequencies of
// all grams across all learned texts
func (gramCollection *GramCollection) updateFrequency(gramIndex int) {
	gramCollection.Frequencies[gramIndex] = gramCollection.Frequencies[gramIndex] + 1
	gramCollection.TotalFrequencies += 1
}

// getIndex fetches the index of a particular gram within the array of grams. If the gram is not found, -1 is returned
func (gramCollection *GramCollection) getIndex(newNgram []string) int {

	for gramIndex, gram := range gramCollection.Grams {

		match := true
		for wordIndex := range gram {
			if gram[wordIndex] != newNgram[wordIndex] {
				match = false
				break
			}
		}

		if match {
			return gramIndex
		}
	}

	return -1
}

// getWeightedRandomNGram returns a random gram from the array of grams, taking the gram's frequency into account.
// First, the gram indices are shuffled to ensure that the grams are iterated in a random order. Next, a random number R
// between 1 and the total frequency count of all grams is generated. Next, the grams are iterated in a random order,
// and for each gram, the corresponding frequency is subtracted from R. When the value of R falls to zero or less, the
// current gram is returned.
func (grams *GramCollection) getWeightedRandomNGram() ([]string, error) {

	grams.RW.Lock()
	defer grams.RW.Unlock()

	if len(grams.Grams) == 0 {
		return []string{}, errors.New("No grams to fetch randomly")
	}

	// if we only have a single gram, return it
	if len(grams.Grams) == 1 {
		return grams.Grams[0], nil
	}

	grams.Shuffle()

	randomIndex := rand.Intn(grams.TotalFrequencies + 1)

	for _, gramIndex := range grams.Indices {

		randomIndex -= grams.Frequencies[gramIndex]

		if randomIndex <= 0 {
			return grams.Grams[gramIndex], nil
		}
	}

	return []string{}, errors.New("Unable to fetch a random n gram")
}

// BuildRandomText returns a random string of text based on the grams learned from the learned texts. First, a random
// gram is selected as the starting point. Next, a subsequent gram is determined, and the last element of the random
// gram is appended to the starting point. This process repeats until no subsequent gram can be determined.
func (grams *GramCollection) BuildRandomText(maxWords, gramSize int) (string, error) {

	startPoint, err := grams.getWeightedRandomNGram()

	if err != nil {
		return "", err
	}

	complete := []string{}

	complete = append(complete, startPoint...)

	nextGram, err := grams.getNext(startPoint, gramSize)

	if err != nil {
		return strings.Join(complete, " "), nil
	}

	for len(nextGram) > 0 {
		nextElement := nextGram[len(nextGram)-1]

		complete = append(complete, nextElement)

		// for unigrams, we need to consider a maximum length
		if maxWords > 0 && len(complete) >= maxWords {
			break
		}

		nextGram, err = grams.getNext(nextGram, gramSize)

		if err != nil {
			break
		}
	}

	return strings.Join(complete, " "), nil
}

// getNext returns a gram from a set of grams whose first two words match the last two words of currentNGram. With the
// set of grams determined, a single gram is randomly selected, taking the gram frequency into account.
func (grams *GramCollection) getNext(currentNGram []string, gramSize int) ([]string, error) {

	grams.RW.RLock()
	defer grams.RW.RUnlock()

	results := NewCollection()

	for _, gramIndex := range grams.Indices {
		gram := grams.Grams[gramIndex]

		matchingGramFound := true

		for gramElementIndex := 0; gramElementIndex < gramSize-1; gramElementIndex++ {

			// check that the last two elements of the previous n gram match the first two elements of the current ngram
			if currentNGram[gramElementIndex+1] != gram[gramElementIndex] {
				matchingGramFound = false
				break
			}
		}

		if matchingGramFound {
			results.Grams = append(results.Grams, gram)
			results.Frequencies = append(results.Frequencies, grams.Frequencies[gramIndex])
			results.Indices = append(results.Indices, len(results.Indices))
			results.TotalFrequencies += grams.Frequencies[gramIndex]
		}
	}

	randomGram, err := results.getWeightedRandomNGram()

	if err != nil {
		return []string{}, err
	}

	return randomGram, nil

}

// Shuffle shuffles the array of indices for the gram collection. This array of indices is used to determine which index
// in the gram collection should be read from. Hence, if we shuffle the array of indices, we effectively read from the
// gram collection in a random order.
func (grams *GramCollection) Shuffle() {

	rand.Seed(time.Now().UnixNano())

	// randomize the gram indices
	rand.Shuffle(len(grams.Indices), func(i, j int) {
		grams.Indices[i], grams.Indices[j] = grams.Indices[j], grams.Indices[i]
	})

}

func (gramCollection *GramCollection) AddGram(newNgram []string) {

	gramCollection.RW.Lock()
	defer gramCollection.RW.Unlock()

	gramIndex := gramCollection.getIndex(newNgram)

	if gramIndex > -1 {
		gramCollection.updateFrequency(gramIndex)
	} else {
		gramCollection.addNewGram(newNgram)
	}
}
