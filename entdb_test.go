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

func GenerateEntVideos(edb *EntDB) []*EntVideo {
	video1 := NewEntVideo(edb)
	video1.Id = uint(123456)
	video1.Title = "title number 1"
	video1.Slug = "title-number-1"
	video1.Tags = []*EntKeyword{
		{Phrase: "tag 1", Type: EntKeywordTag},
		{Phrase: "tag 2", Type: EntKeywordTag},
	}
	video1.Models = []*EntKeyword{
		{Phrase: "model 1", Type: EntKeywordModel},
		{Phrase: "model 2", Type: EntKeywordModel},
	}
	video1.Keywords = []*EntKeyword{
		{Phrase: "aaa bbb ccc", Type: EntKeywordKeyword},
		{Phrase: "bbb ccc ddd", Type: EntKeywordKeyword},
	}
	video2 := &EntVideo{
		Id:    uint(123457),
		Title: "title number 2",
		Slug:  "title-number-2",
		Tags: []*EntKeyword{
			{Phrase: "tag 2", Type: EntKeywordTag},
			{Phrase: "tag 3", Type: EntKeywordTag},
		},
		Models: []*EntKeyword{
			{Phrase: "model 3", Type: EntKeywordModel},
		},
		Keywords: []*EntKeyword{
			{Phrase: "ccc ddd fff", Type: EntKeywordKeyword},
			{Phrase: "bbb ccc ddd", Type: EntKeywordKeyword},
		},
	}
	video3 := &EntVideo{
		Id:    uint(123458),
		Title: "title number 3",
		Slug:  "title-number-3",
		Tags: []*EntKeyword{
			{Phrase: "tag 3", Type: EntKeywordTag},
			{Phrase: "tag 4", Type: EntKeywordTag},
			{Phrase: "tag 5", Type: EntKeywordTag},
			{Phrase: "tag 6", Type: EntKeywordTag},
		},
		Models: []*EntKeyword{
			{Phrase: "model 4", Type: EntKeywordModel},
			{Phrase: "model 5", Type: EntKeywordModel},
			{Phrase: "model 6", Type: EntKeywordModel},
		},
		Keywords: []*EntKeyword{
			{Phrase: "#Aaaa bbbb cccc", Type: EntKeywordKeyword},
			{Phrase: "#CCCC dddd ssss", Type: EntKeywordKeyword},
		},
	}

	video4 := &EntVideo{
		Id:    uint(123459),
		Title: "title number 4 (no keywords)",
		Slug:  "title-number-4-no-keywords",
		Tags: []*EntKeyword{
			{Phrase: "tag 3", Type: EntKeywordTag},
			{Phrase: "tag 4", Type: EntKeywordTag},
			{Phrase: "tag 5", Type: EntKeywordTag},
			{Phrase: "tag 6", Type: EntKeywordTag},
		},
		Models: []*EntKeyword{
			{Phrase: "model 4", Type: EntKeywordModel},
			{Phrase: "model 5", Type: EntKeywordModel},
			{Phrase: "model 6", Type: EntKeywordModel},
		},
		Keywords: []*EntKeyword{},
	}

	video5 := &EntVideo{
		Id:    uint(123460),
		Title: "TiTlE number #5 (\"no keywords'!!!!!)",
		Slug:  "title-number-5-no-keywords",
		Tags: []*EntKeyword{
			{Phrase: "tag 3", Type: EntKeywordTag},
			{Phrase: "tag 4", Type: EntKeywordTag},
			{Phrase: "tag 5", Type: EntKeywordTag},
			{Phrase: "tag 6", Type: EntKeywordTag},
		},
		Models: []*EntKeyword{
			{Phrase: "model 4", Type: EntKeywordModel},
			{Phrase: "model 5", Type: EntKeywordModel},
			{Phrase: "model 6", Type: EntKeywordModel},
		},
		Keywords: []*EntKeyword{},
	}

	video6 := &EntVideo{
		Id:    uint(123461),
		Title: "TiTlE with a stop word number #6 (\"no keywords'!!!!!)",
		Slug:  "title-number-6-no-keywords",
		Tags: []*EntKeyword{
			{Phrase: "tag 3", Type: EntKeywordTag},
			{Phrase: "tag 4", Type: EntKeywordTag},
			{Phrase: "tag 5", Type: EntKeywordTag},
			{Phrase: "tag 6", Type: EntKeywordTag},
		},
		Models: []*EntKeyword{
			{Phrase: "model 4", Type: EntKeywordModel},
			{Phrase: "model 5", Type: EntKeywordModel},
			{Phrase: "model 6", Type: EntKeywordModel},
		},
		Keywords: []*EntKeyword{},
	}

	return []*EntVideo{video1, video2, video3, video4, video5, video6}
}

func TestEntDBAdd(t *testing.T) {
	entdb := NewEntDB("/tmp")

	videos := GenerateEntVideos(entdb)

	for _, video := range videos {
		entdb.Add(video)
	}

	Expected := 6
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

	Expected = 11
	Got = len(entdb.Keywords)
	if Got != Expected {
		t.Errorf("test add video to DB keywords count failed: got %v, wanted %v", Got, Expected)
		for kw, ev := range entdb.Keywords {
			t.Errorf("kw %s -> %d\n", kw, ev.Id)
		}
	}

	Expected = 16
	Got = len(entdb.TwoGrams)
	if Got != Expected {
		t.Errorf("test add video to DB 2-grams count failed: got %v, wanted %v", Got, Expected)
	}

	Expected = 16
	Got = len(entdb.ThreeGrams)
	if Got != Expected {
		t.Errorf("test add video to DB 3-grams count failed: got %v, wanted %v", Got, Expected)
	}
}

func TestEntDBGetVideo(t *testing.T) {
	entdb := NewEntDB("/tmp")

	videos := GenerateEntVideos(entdb)

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

	Expected = entdb.Items[2]
	Got, err = entdb.GetVideoById(uint(123458))
	if err != nil {
		t.Errorf("test get video by id failed: %v", err)
	}
	if Got != Expected {
		t.Errorf("test get video by id failed: got %v, wanted %v", Got, Expected)
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
