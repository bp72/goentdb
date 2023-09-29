package goentdb

import (
	"math/rand"
	"testing"
)

func TestEntVideoGetThumb(t *testing.T) {
	entdb := NewEntDB("/tmp")
	entdb.ThumbBaseUrl = "https://cdn.domain.com/pics"

	videos := GenerateEntVideos(entdb)
	Expected := []string{
		"https://cdn.domain.com/pics/12/34/123456_0.webp",
		"https://cdn.domain.com/pics/12/34/123457_0.webp",
		"https://cdn.domain.com/pics/12/34/123458_0.webp",
		"https://cdn.domain.com/pics/12/34/123459_0.webp",
		"https://cdn.domain.com/pics/12/34/123460_0.webp",
		"https://cdn.domain.com/pics/12/34/123461_0.webp",
	}
	for pos, video := range videos {
		entdb.Add(video)
		Got := video.GetThumb()
		if Got != Expected[pos] {
			t.Errorf("TestEntDBThumbUrls.Test video thumb. Failed: got %v, wanted %v", Got, Expected[pos])
		}
	}

}

func TestEntVideoGetSlug(t *testing.T) {
	entdb := NewEntDB("/tmp")

	videos := GenerateEntVideos(entdb)
	Expected := []string{
		"title-number-1",
		"title-number-2",
		"title-number-3",
		"title-number-4-no-keywords",
		"title-number-5-no-keywords",
		"title-with-a-stop-word-number-6-no-keywords",
	}
	for pos, video := range videos {
		entdb.Add(video)
		Got := video.GetSlug()
		if Got != Expected[pos] {
			t.Errorf("TestEntDBThumbUrls.Test video thumb. Failed: got %v, wanted %v", Got, Expected[pos])
		}
	}

}

func TestEntVideoGetRandomKeyword(t *testing.T) {
	entdb := NewEntDB("/tmp")
	entdb.ThumbBaseUrl = "https://cdn.domain.com/pics"

	videos := GenerateEntVideos(entdb)
	Expected := []string{
		"bbb ccc ddd",
		"bbb ccc ddd",
		"#CCCC dddd ssss",
		"not-found",
		"not-found",
		"not-found",
	}
	for pos, video := range videos {
		entdb.Add(video)
		rand.Seed(1)
		Got := video.GetRandomKeyword()
		if Got.Phrase != Expected[pos] {
			t.Errorf("TestEntDBThumbUrls.Test video thumb. Failed: got %v, wanted %v", Got, Expected[pos])
		}
	}
}

func TestEntVideoGetTokens(t *testing.T) {
	entdb := NewEntDB("/tmp")

	videos := GenerateEntVideos(entdb)
	Expected := [][]string{
		{"title", "number", "1"},
		{"title", "number", "2"},
		{"title", "number", "3"},
		{"title", "number", "4", "no", "keywords"},
		{"title", "number", "5", "no", "keywords"},
		{"title", "with", "a", "stop", "word", "number", "6", "no", "keywords"},
	}
	for pos, video := range videos {
		Got := video.GetTitleTokens(false)
		if len(Got) != len(Expected[pos]) {
			t.Errorf("TestEntVideoGetTokens.Test video thumb. Failed: got %v, wanted %v", Got, Expected[pos])
		}
	}

	Expected = [][]string{
		{"title", "number", "1"},
		{"title", "number", "2"},
		{"title", "number", "3"},
		{"title", "number", "4", "keywords"},
		{"title", "number", "5", "keywords"},
		{"title", "stop", "word", "number", "6", "keywords"},
	}
	for pos, video := range videos {
		Got := video.GetTitleTokens(true)
		if len(Got) != len(Expected[pos]) {
			t.Errorf("TestEntVideoGetTokens.Test grams len. Failed: got %v, wanted %v", Got, Expected[pos])
		}
		for i, got := range Got {
			if got != Expected[pos][i] {
				t.Errorf("TestEntVideoGetTokens.Test gram item. Failed: got %v, wanted %v", got, Expected[pos][i])
			}
		}
	}
}

func TestEntVideoGetNGrams(t *testing.T) {
	entdb := NewEntDB("/tmp")

	videos := GenerateEntVideos(entdb)
	Expected2Grams := [][]string{
		{"title number", "number 1"},
		{"title number", "number 2"},
		{"title number", "number 3"},
		{"title number", "number 4", "4 no", "no keywords"},
		{"title number", "number 5", "5 no", "no keywords"},
		{"title with", "with a", "a stop", "stop word", "word number", "number 6", "6 no", "no keywords"},
	}
	Expected3Grams := [][]string{
		{"title number 1"},
		{"title number 2"},
		{"title number 3"},
		{"title number 4", "number 4 no", "4 no keywords"},
		{"title number 5", "number 5 no", "5 no keywords"},
		{"title with a", "with a stop", "a stop word", "stop word number", "word number 6", "number 6 no", "6 no keywords"},
	}
	for pos, video := range videos {
		Got2Grams, Got3Grams := video.GetNGrams(false)
		if len(Got2Grams) != len(Expected2Grams[pos]) {
			t.Errorf("TestEntVideoGetNGrams.Test 2gram len. Failed: got %v, wanted %v", Got2Grams, Expected2Grams[pos])
		}
		for i, got := range Got2Grams {
			if got != Expected2Grams[pos][i] {
				t.Errorf("TestEntVideoGetNGrams.Test 2gram item. Failed: got %v, wanted %v", Got2Grams, Expected2Grams[pos][i])
			}
		}
		if len(Got3Grams) != len(Expected3Grams[pos]) {
			t.Errorf("TestEntVideoGetNGrams.Test 3gram len. Failed: got %v, wanted %v", Got3Grams, Expected3Grams[pos])
		}
		for i, got := range Got3Grams {
			if got != Expected3Grams[pos][i] {
				t.Errorf("TestEntVideoGetNGrams.Test 3gram item. Failed: got %v, wanted %v", Got3Grams, Expected3Grams[pos][i])
			}
		}
	}
}

func TestEntVideoGetNGramsWithoutStopWords(t *testing.T) {
	entdb := NewEntDB("/tmp")

	videos := GenerateEntVideos(entdb)
	Expected2Grams := [][]string{
		{"title number", "number 1"},
		{"title number", "number 2"},
		{"title number", "number 3"},
		{"title number", "number 4", "4 keywords"},
		{"title number", "number 5", "5 keywords"},
		{"title stop", "stop word", "word number", "number 6", "6 keywords"},
	}
	Expected3Grams := [][]string{
		{"title number 1"},
		{"title number 2"},
		{"title number 3"},
		{"title number 4", "number 4 keywords"},
		{"title number 5", "number 5 keywords"},
		{"title stop word", "stop word number", "word number 6", "number 6 keywords"},
	}
	for pos, video := range videos {
		Got2Grams, Got3Grams := video.GetNGrams(true)
		if len(Got2Grams) != len(Expected2Grams[pos]) {
			t.Errorf("TestEntVideoGetNGrams.Test 2gram len. Failed: got %v, wanted %v", Got2Grams, Expected2Grams[pos])
		}
		for i, got := range Got2Grams {
			if got != Expected2Grams[pos][i] {
				t.Errorf("TestEntVideoGetNGrams.Test 2gram item. Failed: got %v, wanted %v", Got2Grams, Expected2Grams[pos][i])
			}
		}
		if len(Got3Grams) != len(Expected3Grams[pos]) {
			t.Errorf("TestEntVideoGetNGrams.Test 3gram len. Failed: got %v, wanted %v", Got3Grams, Expected3Grams[pos])
		}
		for i, got := range Got3Grams {
			if got != Expected3Grams[pos][i] {
				t.Errorf("TestEntVideoGetNGrams.Test 3gram item. Failed: got %v, wanted %v", Got3Grams, Expected3Grams[pos][i])
			}
		}
	}
}
