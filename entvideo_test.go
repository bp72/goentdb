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
	}
	for pos, video := range videos {
		entdb.Add(video)
		Got := video.GetThumb()
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
