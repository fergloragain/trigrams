package gram

import (
	"fmt"
	"strings"
	"testing"
)

func TestShuffle(t *testing.T) {
	grams := NewCollection()

	grams.Indices = []int{0, 1, 2, 4, 5}

	grams.Shuffle()

	for _, v := range grams.Indices {
		fmt.Println(v)
	}
}

func TestGetIndex(t *testing.T) {
	grams := NewCollection()

	tt := []struct {
		Grams         [][]string
		Search        []string
		ExpectedIndex int
	}{
		{
			Grams: [][]string{
				{"dog", "sit", "bark"},
				{"cat", "jump", "meow"},
				{"owl", "fly", "hoot"},
				{"fish", "swim", "blub"},
			},
			Search:        []string{"cat", "jump", "meow"},
			ExpectedIndex: 1,
		},
		{
			Grams: [][]string{
				{"dog", "sit", "bark"},
				{"cat", "jump", "meow"},
				{"owl", "fly", "hoot"},
				{"fish", "swim", "blub"},
			},
			Search:        []string{"owl", "jump", "blub"},
			ExpectedIndex: -1,
		},
		{
			Grams: [][]string{
				{"dog", "sit", "bark"},
				{"cat", "jump", "meow"},
				{"owl", "fly", "hoot"},
				{"fish", "swim", "blub"},
			},
			Search:        []string{"fish", "swim", "blub"},
			ExpectedIndex: 3,
		},
		{
			Grams: [][]string{
				{"fish", "swim", "blub"},
			},
			Search:        []string{"fish", "swim", "blub"},
			ExpectedIndex: 0,
		},
	}

	for _, v := range tt {
		grams.Grams = v.Grams

		index := grams.getIndex(v.Search)

		if index != v.ExpectedIndex {
			t.Fail()
		}
	}

}

func TestAddNewGram(t *testing.T) {

	tt := []struct {
		GramsToAdd       [][]string
		Grams            [][]string
		Frequencies      []int
		Indices          []int
		TotalFrequencies int
	}{
		{
			GramsToAdd:       [][]string{{"fish", "swim", "blub"}},
			Grams:            [][]string{{"fish", "swim", "blub"}},
			Frequencies:      []int{1},
			Indices:          []int{0},
			TotalFrequencies: 1,
		},
		{
			GramsToAdd:       [][]string{{"fish", "swim", "blub"}, {"shark", "swim", "blub"}},
			Grams:            [][]string{{"fish", "swim", "blub"}, {"shark", "swim", "blub"}},
			Frequencies:      []int{1, 1},
			Indices:          []int{0, 1},
			TotalFrequencies: 2,
		},
	}

	for _, v := range tt {
		grams := NewCollection()

		for _, g := range v.GramsToAdd {
			grams.addNewGram(g)
		}

		if len(v.Grams) != len(grams.Grams) {
			t.Error("len(v.Grams) != len(grams.Grams)")
		}

		for i := 0; i < len(v.GramsToAdd); i++ {
			g1 := v.GramsToAdd[i]
			g2 := grams.Grams[i]

			for j := range g1 {
				if g1[j] != g2[j] {
					t.Error("g1 != g2")
				}
			}
		}

		if grams.TotalFrequencies != v.TotalFrequencies {
			t.Error("grams.TotalFrequencies != v.TotalFrequencies")
		}

		for i := range grams.Frequencies {
			if v.Frequencies[i] != grams.Frequencies[i] {
				t.Error("v.Frequencies[i] != grams.Frequencies[i]")
			}
		}

		for i := range grams.Indices {
			if v.Indices[i] != grams.Indices[i] {
				t.Error("v.IndicesToUpdate[i] != grams.IndicesToUpdate[i]")
			}
		}
	}

}

func TestUpdateFrequency(t *testing.T) {

	tt := []struct {
		StartingFrequencies []int
		IndicesToUpdate     []int
		ExpectedFrequencies []int
		TotalFrequencies    int
	}{
		{
			StartingFrequencies: []int{0, 0, 0, 0, 0},
			IndicesToUpdate:     []int{0, 0, 1, 2, 3, 4, 4, 4, 0, 0},
			ExpectedFrequencies: []int{4, 1, 1, 1, 3},
			TotalFrequencies:    10,
		},
		{
			StartingFrequencies: []int{0, 0, 0, 0, 0},
			IndicesToUpdate:     []int{0, 0, 0, 0, 0, 0, 0, 0},
			ExpectedFrequencies: []int{8, 0, 0, 0, 0},
			TotalFrequencies:    8,
		},
	}

	for _, v := range tt {

		grams := NewCollection()

		grams.Frequencies = v.StartingFrequencies

		for _, v := range v.IndicesToUpdate {
			grams.updateFrequency(v)
		}

		for i := range v.ExpectedFrequencies {
			if grams.Frequencies[i] != v.ExpectedFrequencies[i] {
				t.Error("Mismatch of frequencies")
			}
		}

		if v.TotalFrequencies != grams.TotalFrequencies {
			t.Error("Mismatch of total frequencies")
		}

	}

}

func TestGetWeightedRandomNGram(t *testing.T) {

	tt := []struct {
		Grams            [][]string
		Frequencies      []int
		TotalFrequencies int
		Indices          []int
	}{
		{
			Grams: [][]string{
				[]string{
					"abc", "def", "ghi",
				},
				[]string{
					"xyz", "123", "mnb",
				},
			},
			Frequencies: []int{
				90,
				10,
			},
			TotalFrequencies: 100,
			Indices: []int{
				0,
				1,
			},
		},
		{
			Grams: [][]string{
				[]string{
					"1", "2", "3",
				},
				[]string{
					"4", "5", "6",
				},
				[]string{
					"7", "8", "9",
				},
			},
			Frequencies: []int{
				30,
				30,
				30,
			},
			TotalFrequencies: 90,
			Indices: []int{
				0,
				1,
				2,
			},
		},
		{
			Grams: [][]string{
				[]string{
					"1", "2", "3",
				},
				[]string{
					"4", "5", "6",
				},
				[]string{
					"7", "8", "9",
				},
			},
			Frequencies: []int{
				10,
				10,
				70,
			},
			TotalFrequencies: 90,
			Indices: []int{
				0,
				1,
				2,
			},
		},
	}

	for _, v := range tt {

		grams := NewCollection()

		grams.Grams = v.Grams
		grams.Indices = v.Indices
		grams.Frequencies = v.Frequencies
		grams.TotalFrequencies = v.TotalFrequencies

		stats := struct {
			Token     [][]string
			Frequency []int
		}{
			Token:     [][]string{},
			Frequency: []int{},
		}

		for i := 0; i < 1000; i++ {
			randomNGram, err := grams.getWeightedRandomNGram()

			if err != nil {
				t.Error(err.Error())
			}

			exists := false

			for wordIndex, wordArray := range stats.Token {
				found := true

				for h, word := range wordArray {
					if word != randomNGram[h] {
						found = false
					}

				}

				if found {
					exists = true
					stats.Frequency[wordIndex] = stats.Frequency[wordIndex] + 1
					break
				}
			}

			if !exists {
				stats.Token = append(stats.Token, randomNGram)
				stats.Frequency = append(stats.Frequency, 1)
			}

		}

		// first pass should reflect a weighting of 9:1 for [abc def ghi]:[xyz 123 mnb]
		// second pass should reflect a weighting of 1:1:1 for [1 2 3]:[4 5 6]:[7 8 9]
		// third pass should reflect a weighting of 7:1:1 for [7 8 9]:[1 2 3]:[4 5 6]
		// e.g.:
		//[[abc def ghi] [xyz 123 mnb]]
		//[917 83]
		//[[7 8 9] [1 2 3] [4 5 6]]
		//[327 327 346]
		//[[7 8 9] [1 2 3] [4 5 6]]
		//[783 113 104]

		fmt.Println(stats.Token)
		fmt.Println(stats.Frequency)

	}

}

func TestGetWeightedRandomNGram_BadData(t *testing.T) {

	tt := []struct {
		Grams            [][]string
		Frequencies      []int
		TotalFrequencies int
		Indices          []int
		Expected         []string
	}{
		{
			Grams:            [][]string{},
			Frequencies:      []int{},
			TotalFrequencies: 0,
			Indices:          []int{},
			Expected:         []string{},
		},
		{
			Grams: [][]string{
				[]string{"1", "2", "3"},
			},
			Frequencies: []int{
				1,
			},
			TotalFrequencies: 1,
			Indices: []int{
				0,
			},
			Expected: []string{"1", "2", "3"},
		},
	}

	for _, v := range tt {

		grams := NewCollection()

		grams.Grams = v.Grams
		grams.Indices = v.Indices
		grams.Frequencies = v.Frequencies
		grams.TotalFrequencies = v.TotalFrequencies

		randomNGram, _ := grams.getWeightedRandomNGram()

		if len(randomNGram) != len(v.Expected) {
			t.Fail()
		}
	}
}

func TestBuildRandomText(t *testing.T) {

	tt := []struct {
		Grams          [][]string
		Frequencies    []int
		Indices        []int
		TotalFrequency int
		ResultSuperset string
		GramSize       int
		Maxwords       int
	}{
		{
			Grams: [][]string{
				[]string{
					"this", "is", "a",
				},
				[]string{
					"is", "a", "sample",
				},
				[]string{
					"a", "sample", "text",
				},
				[]string{
					"sample", "text", "blob",
				},
			},
			Frequencies: []int{
				1, 1, 1, 1,
			},
			Indices: []int{
				0, 1, 2, 3,
			},
			TotalFrequency: 4,
			ResultSuperset: "this is a sample text blob",
			GramSize:       3,
			Maxwords:       100,
		},
		{
			Grams:          [][]string{},
			Frequencies:    []int{},
			Indices:        []int{},
			TotalFrequency: 0,
			ResultSuperset: "",
			GramSize:       3,
			Maxwords:       100,
		},
		{
			Grams: [][]string{
				[]string{"a", "b", "c"},
				[]string{"d", "e", "f"},
			},
			Frequencies:    []int{},
			Indices:        []int{},
			TotalFrequency: 0,
			ResultSuperset: "",
			GramSize:       3,
			Maxwords:       100,
		},
		{
			Grams: [][]string{
				[]string{
					"this", "is", "a",
				},
				[]string{
					"is", "a", "sample",
				},
				[]string{
					"is", "this", "is",
				},
			},
			Frequencies: []int{
				1, 1, 1,
			},
			Indices: []int{
				0, 1, 2,
			},
			TotalFrequency: 3,
			ResultSuperset: "is this is this is a sample",
			GramSize:       3,
			Maxwords:       100,
		},
		{
			Grams: [][]string{
				[]string{
					"this", "is", "a",
				},
				[]string{
					"is", "a", "sample",
				},
				[]string{
					"is", "this", "is",
				},
			},
			Frequencies: []int{
				1, 1, 1,
			},
			Indices: []int{
				0,
			},
			TotalFrequency: 3,
			ResultSuperset: "is this is this is a sample",
			GramSize:       3,
			Maxwords:       100,
		},
		{
			Grams: [][]string{
				[]string{
					"this", "that", "those",
				},
				[]string{
					"that", "those", "this",
				},
				[]string{
					"those", "this", "that",
				},
			},
			Frequencies: []int{
				1, 1, 1,
			},
			Indices: []int{
				0, 1, 2,
			},
			TotalFrequency: 3,
			ResultSuperset: "this that those this that those",
			GramSize:       3,
			Maxwords:       1,
		},
		{
			Grams: [][]string{
				[]string{
					"this", "that", "those",
				},
				[]string{
					"xyz", "abc", "def",
				},
				[]string{
					"1", "2", "3",
				},
			},
			Frequencies: []int{
				-1,
			},
			Indices: []int{
				0,
			},
			TotalFrequency: 3000,
			ResultSuperset: "this that those xyz abc def 1 2 3",
			GramSize:       3,
			Maxwords:       1,
		},
	}

	for _, v := range tt {
		grams := NewCollection()

		grams.Grams = v.Grams
		grams.Frequencies = v.Frequencies
		grams.TotalFrequencies = v.TotalFrequency
		grams.Indices = v.Indices

		randomString, _ := grams.BuildRandomText(v.Maxwords, v.GramSize)

		if !strings.Contains(v.ResultSuperset, randomString) {
			t.Error(fmt.Sprintf("%s does not contain %s", v.ResultSuperset, randomString))
		}
	}

}

func TestGetNext(t *testing.T) {

	tt := []struct {
		Grams          [][]string
		Frequencies    []int
		Indices        []int
		TotalFrequency int
		GramSize       int
		CurrentGram    []string
		ExpectedGram   []string
	}{
		{
			Grams: [][]string{
				[]string{
					"this", "is", "a",
				},
				[]string{
					"is", "a", "sample",
				},
				[]string{
					"a", "sample", "text",
				},
				[]string{
					"sample", "text", "blob",
				},
			},
			Frequencies: []int{
				1, 1, 1, 1,
			},
			Indices: []int{
				0, 1, 2, 3,
			},
			TotalFrequency: 4,
			GramSize:       3,
			CurrentGram: []string{
				"this", "is", "a",
			},
			ExpectedGram: []string{
				"is", "a", "sample",
			},
		},
		{
			Grams: [][]string{
				[]string{
					"this", "is", "a",
				},
				[]string{
					"is", "a", "sample",
				},
				[]string{
					"is", "a", "test",
				},
			},
			Frequencies: []int{
				1, 1, 1, 1,
			},
			Indices:        []int{},
			TotalFrequency: 4,
			GramSize:       3,
			CurrentGram: []string{
				"this", "is", "a",
			},
			ExpectedGram: []string{},
		},
	}

	for _, v := range tt {
		grams := NewCollection()

		grams.Grams = v.Grams
		grams.Frequencies = v.Frequencies
		grams.TotalFrequencies = v.TotalFrequency
		grams.Indices = v.Indices

		nextGram, _ := grams.getNext(v.CurrentGram, v.GramSize)

		for i := range nextGram {
			if nextGram[i] != v.ExpectedGram[i] {
				t.Fail()
			}
		}
	}

}

func TestAddGram(t *testing.T) {

	tt := []struct {
		Grams                  [][]string
		Frequencies            []int
		ExpectedFrequencies    []int
		Indices                []int
		TotalFrequency         int
		ExpectedTotalFrequency int
		GramSize               int
		CurrentGram            []string
		GramCollection         *GramCollection
	}{
		{
			Grams: [][]string{
				[]string{
					"this", "is", "a",
				},
				[]string{
					"is", "a", "sample",
				},
				[]string{
					"a", "sample", "text",
				},
				[]string{
					"sample", "text", "blob",
				},
			},
			Frequencies: []int{
				1, 1, 1, 1,
			},
			ExpectedFrequencies: []int{
				2, 1, 1, 1,
			},
			Indices: []int{
				0, 1, 2, 3,
			},
			TotalFrequency:         4,
			ExpectedTotalFrequency: 5,

			GramSize: 3,
			CurrentGram: []string{
				"this", "is", "a",
			},
		},

		{
			Grams: [][]string{
				[]string{
					"this", "is", "a",
				},
				[]string{
					"is", "a", "sample",
				},
				[]string{
					"a", "sample", "text",
				},
				[]string{
					"sample", "text", "blob",
				},
			},
			Frequencies: []int{
				1, 1, 1, 1,
			},
			ExpectedFrequencies: []int{
				1, 1, 1, 1, 1,
			},
			Indices: []int{
				0, 1, 2, 3,
			},
			TotalFrequency:         4,
			ExpectedTotalFrequency: 5,
			GramSize:               3,
			CurrentGram: []string{
				"this", "isn't", "a",
			},
		},
	}

	for _, v := range tt {
		grams := NewCollection()

		grams.Grams = v.Grams
		grams.Frequencies = v.Frequencies
		grams.TotalFrequencies = v.TotalFrequency
		grams.Indices = v.Indices

		grams.AddGram(v.CurrentGram)

		for i := range grams.Frequencies {
			if grams.Frequencies[i] != v.ExpectedFrequencies[i] {
				t.Fail()
			}
		}
	}

}
