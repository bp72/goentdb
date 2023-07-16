package goentdb

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestEntDBCreate(t *testing.T) {
	entdb := NewEntDB("/tmp")

	Expected := "/tmp"
	Got := entdb.StoragePath
	if Got != Expected {
		t.Errorf("test path failed: got %v, wanted %v", Got, Expected)
	}
}

func GenerateEntVideos() []*EntVideo {
	video1 := &EntVideo{
		Id:    uint(1),
		Title: "title number 1",
		Slug:  "title-number-1",
		Tags: []*EntKeyword{
			&EntKeyword{Phrase: "tag 1", Type: EntKeywordTag},
			&EntKeyword{Phrase: "tag 2", Type: EntKeywordTag},
		},
		Models: []*EntKeyword{
			&EntKeyword{Phrase: "model 1", Type: EntKeywordModel},
			&EntKeyword{Phrase: "model 2", Type: EntKeywordModel},
		},
		Keywords: []*EntKeyword{
			&EntKeyword{Phrase: "aaa bbb ccc", Type: EntKeywordKeyword},
			&EntKeyword{Phrase: "bbb ccc ddd", Type: EntKeywordKeyword},
		},
	}
	video2 := &EntVideo{
		Id:    uint(1),
		Title: "title number 2",
		Slug:  "title-number-2",
		Tags: []*EntKeyword{
			&EntKeyword{Phrase: "tag 2", Type: EntKeywordTag},
			&EntKeyword{Phrase: "tag 3", Type: EntKeywordTag},
		},
		Models: []*EntKeyword{
			&EntKeyword{Phrase: "model 3", Type: EntKeywordModel},
		},
		Keywords: []*EntKeyword{
			&EntKeyword{Phrase: "ccc ddd fff", Type: EntKeywordKeyword},
			&EntKeyword{Phrase: "bbb ccc ddd", Type: EntKeywordKeyword},
		},
	}
	video3 := &EntVideo{
		Id:    uint(1),
		Title: "title number 3",
		Slug:  "title-number-3",
		Tags: []*EntKeyword{
			&EntKeyword{Phrase: "tag 3", Type: EntKeywordTag},
			&EntKeyword{Phrase: "tag 4", Type: EntKeywordTag},
			&EntKeyword{Phrase: "tag 5", Type: EntKeywordTag},
			&EntKeyword{Phrase: "tag 6", Type: EntKeywordTag},
		},
		Models: []*EntKeyword{
			&EntKeyword{Phrase: "model 4", Type: EntKeywordModel},
			&EntKeyword{Phrase: "model 5", Type: EntKeywordModel},
			&EntKeyword{Phrase: "model 6", Type: EntKeywordModel},
		},
		Keywords: []*EntKeyword{
			&EntKeyword{Phrase: "#Aaaa bbbb cccc", Type: EntKeywordKeyword},
			&EntKeyword{Phrase: "#CCCC dddd ssss", Type: EntKeywordKeyword},
		},
	}

	return []*EntVideo{video1, video2, video3}
}

func TestEntDBAdd(t *testing.T) {
	entdb := NewEntDB("/tmp")

	videos := GenerateEntVideos()

	for _, video := range videos {
		entdb.Add(video)
	}

	Expected := 3
	Got := len(entdb.Items)
	if Got != Expected {
		t.Errorf("test add video to DB failed: got %v, wanted %v", Got, Expected)
	}

	Expected = 6
	Got = len(entdb.Tags)
	if Got != Expected {
		t.Errorf("test add video to DB tags count failed: got %v, wanted %v", Got, Expected)
	}

	Expected = 6
	Got = len(entdb.Models)
	if Got != Expected {
		t.Errorf("test add video to DB models count failed: got %v, wanted %v", Got, Expected)
	}

	Expected = 8
	Got = len(entdb.Keywords)
	if Got != Expected {
		t.Errorf("test add video to DB keywords count failed: got %v, wanted %v", Got, Expected)
	}
}

func TestEntDBGetVideo(t *testing.T) {
	entdb := NewEntDB("/tmp")

	videos := GenerateEntVideos()

	for _, video := range videos {
		entdb.Add(video)
	}

	Expected := videos[0]
	Got, err := entdb.GetVideoByMD5(MD5("title-number-1"))
	if err != nil {
		t.Errorf("test get video by md5 failed: %v", err)
	}
	if Got != Expected {
		t.Errorf("test get video by md5 failed: got %v, wanted %v", Got, Expected)
	}

	_, err = entdb.GetVideoByMD5(MD5("not-existing-video"))
	if err == nil {
		t.Errorf("test get video by md5 failed: %v", err)
	}
	ExpectedErr := "EntVideo not found"
	if err.Error() != ExpectedErr {
		t.Errorf("test get error when video not found by md5 failed: got %v, wanted %v", err, ExpectedErr)
	}

	Got, err = entdb.GetVideoByMD5(MD5("aaa-bbb-ccc"))
	if err != nil {
		t.Errorf("test get video by keyword md5 failed: %v", err)
	}
	if Got != Expected {
		t.Errorf("test get video by keyword md5 failed: got %v, wanted %v", Got, Expected)
	}
}

func TestEntDBDumpLoadTags(t *testing.T) {
	entdb := NewEntDB("/tmp")

	if _, err := os.Stat(entdb.GetDictTagsPath()); !errors.Is(err, os.ErrNotExist) {
		e := os.Remove(entdb.GetDictTagsPath())
		if e != nil {
			t.Errorf("Error deleting %s: %v", entdb.GetDictTagsPath(), e)
		}
	}

	for i := 0; i < 10; i++ {
		entdb.AddTag(NewTag(i, fmt.Sprintf("tag-%d", i)))
	}

	Expected := 10
	Got := len(entdb.DictTags)
	if Got != Expected {
		t.Errorf("test tags should be %d got %d", Expected, Got)
	}

	entdb.DumpTags()

	entdb_new := NewEntDB("/tmp")
	entdb_new.LoadTags()

	Got = len(entdb_new.DictTags)
	if Got != Expected {
		t.Errorf("test tags should be %d got %d", Expected, Got)
	}

	for i := 0; i < 10; i++ {
		a := entdb.DictTags[i]
		b := entdb_new.DictTags[i]
		if a.Phrase != b.Phrase || a.Type != b.Type {
			t.Errorf("test tag %d should be %s got %s", i, a.Phrase, b.Phrase)
		}
	}
}

func TestEntDBDumpLoadModels(t *testing.T) {
	entdb := NewEntDB("/tmp")

	if _, err := os.Stat(entdb.GetDictModelsPath()); !errors.Is(err, os.ErrNotExist) {
		e := os.Remove(entdb.GetDictModelsPath())
		if e != nil {
			t.Errorf("Error deleting %s: %v", entdb.GetDictModelsPath(), e)
		}
	}

	for i := 0; i < 10; i++ {
		entdb.AddModel(NewModel(i, fmt.Sprintf("model-%d", i)))
	}

	Expected := 10
	Got := len(entdb.DictModels)
	if Got != Expected {
		t.Errorf("test models should be %d got %d", Expected, Got)
	}

	entdb.DumpModels()

	entdb_new := NewEntDB("/tmp")
	entdb_new.LoadModels()

	Got = len(entdb_new.DictModels)
	if Got != Expected {
		t.Errorf("test models should be %d got %d", Expected, Got)
	}

	for i := 0; i < 10; i++ {
		a := entdb.DictModels[i]
		b := entdb_new.DictModels[i]
		if a.Phrase != b.Phrase || a.Type != b.Type {
			t.Errorf("test model %d should be %s got %s", i, a.Phrase, b.Phrase)
		}
	}
}
