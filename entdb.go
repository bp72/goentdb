package goentdb

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"
)

//Запилить структуру для загрузки данных
//Запилить ссылки с тег -[]*EnvtVideo

type EntDB struct {
	StoragePath  string
	Items        []*EntVideo
	SeoPool      []*EntKeyword
	Tags         map[string][]*EntVideo
	Models       map[string][]*EntVideo
	Search       map[string][]*EntVideo
	Keywords     map[string]*EntVideo
	DictTags     map[int]*EntKeyword
	DictModels   map[int]*EntKeyword
	lock         sync.RWMutex
	Origins      map[Origin]int
	ThumbBaseUrl string
}

func (edb *EntDB) GetDictTagsPath() string {
	return fmt.Sprintf("%s/tags", edb.StoragePath)
}

func (edb *EntDB) GetDictModelsPath() string {
	return fmt.Sprintf("%s/models", edb.StoragePath)
}

func (edb *EntDB) GetDictVideosPath() string {
	return fmt.Sprintf("%s/videos", edb.StoragePath)
}

func (edb *EntDB) GetTagById(id int) (*EntKeyword, error) {
	edb.lock.RLock()
	defer edb.lock.RUnlock()

	if tag, exists := edb.DictTags[id]; exists {
		return tag, nil
	}

	return nil, errors.New(fmt.Sprintf("EntKeyword(Tag) not found: %d", id))
}

func (edb *EntDB) GetModelById(id int) (*EntKeyword, error) {
	edb.lock.RLock()
	defer edb.lock.RUnlock()

	if model, exists := edb.DictModels[id]; exists {
		return model, nil
	}

	return nil, errors.New(fmt.Sprintf("EntKeyword(Model) not found: %d", id))
}

func (edb *EntDB) AddTag(tag *EntKeyword) {
	edb.lock.Lock()
	defer edb.lock.Unlock()

	edb.DictTags[tag.Id] = tag
}

func (edb *EntDB) AddModel(model *EntKeyword) {
	edb.lock.Lock()
	defer edb.lock.Unlock()

	edb.DictModels[model.Id] = model
}

/*
	Add video to the DB
	- Add to slice for future slices and random access
	- Add to tag/model -> []*EntVideo for tag/model slices
	- Add to keyword -> *EntVideos map for original slug md5 and keyword slug md5 access
*/
func (edb *EntDB) Add(video *EntVideo) {
	edb.lock.Lock()
	defer edb.lock.Unlock()

	video.Owner = edb
	edb.Items = append(edb.Items, video)
	for _, tag := range video.Tags {
		edb.Tags[tag.GetSlug()] = append(edb.Tags[tag.GetSlug()], video)
	}
	for _, model := range video.Models {
		edb.Models[model.GetSlug()] = append(edb.Models[model.GetSlug()], video)
	}

	// Add original slug to map for O(1) access for video
	edb.Keywords[video.GetMD5()] = video

	for _, keyword := range video.Keywords {
		edb.Keywords[keyword.GetMD5()] = video
	}

	edb.Origins[video.Origin]++

	Tokens := strings.Split(strings.ToLower(video.Title), " ")
	for _, token := range Tokens {
		if len(token) < 3 {
			continue
		}
		videos, _ := edb.Search[token]
		edb.Search[token] = append(videos, video)
	}

}

func (edb *EntDB) AddVideoFromLoad(evfl *EntVideoForLoad) {
	ev := NewEntVideo(edb)

	ev.Id = evfl.Id
	ev.Title = evfl.Title
	ev.Origin = evfl.Origin
	ev.OriginId = evfl.OriginId
	ev.OriginUrl = evfl.OriginUrl
	ev.Duration = evfl.Duration
	ev.Slug = evfl.Slug
	ev.Source = evfl.Source
	ev.Descr = evfl.Descr
	ev.ModifiedAt = evfl.ModifiedAt
	ev.Keywords = evfl.Keywords
	ev.ThumbUrls = evfl.ThumbUrls
	ev.VideoUrls = evfl.VideoUrls

	for _, tag_id := range evfl.Tags {
		tag, err := edb.GetTagById(tag_id)
		if err != nil {
			log.Fatal(err)
		}
		ev.AddTag(tag)
	}

	for _, model_id := range evfl.Models {
		model, err := edb.GetModelById(model_id)
		if err != nil {
			log.Fatal(err)
		}
		ev.AddModel(model)
	}

	// edb.Items = append(edb.Items, ev)
	edb.Add(ev)
}

/*
	Get Video by original slug md5 or keyword slug md5
*/
func (edb *EntDB) GetVideoByMD5(key string) (*EntVideo, error) {
	edb.lock.RLock()
	defer edb.lock.RUnlock()

	if video, exists := edb.Keywords[key]; exists {
		return video, nil
	}

	return nil, errors.New("EntVideo not found")
}

func (edb *EntDB) GetKeywordsRandomSet(Size int, UseSeoPool bool, Exclude []*EntVideo) []*EntKeyword {
	res := make([]*EntKeyword, Size)

	pos := 0

	for pos < Size {
		ev := edb.Random()
		// TODO: if (ev not in Exclude and ev not in Taken) {
		// Check if EntVideo has any associated EntKeyword
		if len(ev.Keywords) > 0 {
			res[pos] = ev.GetRandomKeyword()
			pos++
		}
		//}
	}

	return res
}

/*
	Get slice of random EntVideos based on query filter
*/
func (edb *EntDB) RandomSetBySearch(Query string, Size int) ([]*EntVideo, int) {
	QueryTokens := strings.Split(strings.ToLower(Query), " ")
	Counter := make(map[*EntVideo]int)

	for _, Token := range QueryTokens {
		if videos, exists := edb.Search[Token]; exists {
			for _, video := range videos {
				Counter[video]++
			}
		}
	}

	type KeyValue struct {
		Video *EntVideo
		Value int
	}

	var SortedSlice []KeyValue

	for k, v := range Counter {
		SortedSlice = append(SortedSlice, KeyValue{k, v})
	}

	sort.Slice(SortedSlice, func(i, j int) bool {
		return SortedSlice[i].Value > SortedSlice[j].Value
	})

	res := make([]*EntVideo, Min(len(SortedSlice), Size))

	for i := 0; i < Min(len(SortedSlice), Size); i++ {
		res[i] = SortedSlice[i].Video
	}

	return res, len(SortedSlice)
}

func (ev *EntDB) RandomSetByModel(ModelSlug string, Size int) ([]*EntVideo, int) {

	models, exists := ev.Models[ModelSlug]

	if !exists {
		return make([]*EntVideo, 0), 0
	}

	var candidate *EntVideo

	if Size >= len(models) {
		return models, len(models)
	}

	seen := make(map[*EntVideo]bool)
	res := make([]*EntVideo, Size)

	for i := 0; i < Size; i++ {
		candidate = models[rand.Intn(len(models))]
		if _, exists := seen[candidate]; !exists {
			res[i] = candidate
			seen[candidate] = true
		}
	}

	return res, len(models)
}

func (ev *EntDB) RandomSetByTag(TagSlug string, Size int) ([]*EntVideo, int) {
	tags, exists := ev.Tags[TagSlug]

	if !exists {
		return make([]*EntVideo, 0), 0
	}

	var candidate *EntVideo

	if Size >= len(tags) {
		return tags, len(tags)
	}

	seen := make(map[*EntVideo]bool)
	res := make([]*EntVideo, Size)

	for i := 0; i < Size; i++ {
		candidate = tags[rand.Intn(len(tags))]
		if _, exists := seen[candidate]; !exists {
			res[i] = candidate
			seen[candidate] = true
		}
	}

	return res, len(tags)
}

func (edb *EntDB) RandomSet(Size int) []*EntVideo {
	res := make([]*EntVideo, Size)

	for i := 0; i < Size; i++ {
		res[i] = edb.Random()
	}

	return res
}

func (edb *EntDB) Random() *EntVideo {
	edb.lock.RLock()
	defer edb.lock.RUnlock()
	return edb.Items[rand.Intn(len(edb.Items))]
}

/*
	Return random set of Keywords of specific size
*/
func (edb *EntDB) RandomKeywordSet(Size int, useSeoPool bool) []*EntKeyword {
	if useSeoPool {
		return edb.RandomKeywordSetFromSeoPool(Size)
	}
	return edb.RandomKeywordSetFromGeneralPool(Size)
}

func (edb *EntDB) RandomKeywordSetFromGeneralPool(Size int) []*EntKeyword {
	res := make([]*EntKeyword, Size)

	for i := 0; i < Size; i++ {
		res[i] = edb.Random().GetRandomKeyword()
	}

	return res
}

func (edb *EntDB) RandomKeywordSetFromSeoPool(Size int) []*EntKeyword {
	return edb.RandomKeywordSetFromGeneralPool(Size)
	// res := make([]*EntKeyword, Size)

	// for i := 0; i < Size; i++ {
	// 	res[i] = edb.SeoPool[rand.Intn(len(edb.SeoPool))]
	// 	// edb.S
	// 	// if kw != nil {
	// 	// 	ans = append(ans, kw)
	// 	// 	self.Shows[kw.SlugMD5]++
	// 	// }
	// }

	// return res
}

func NewEntDB(path string) *EntDB {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	return &EntDB{
		StoragePath: path,
		Tags:        make(map[string][]*EntVideo),
		Models:      make(map[string][]*EntVideo),
		Search:      make(map[string][]*EntVideo),
		Keywords:    make(map[string]*EntVideo),
		DictTags:    make(map[int]*EntKeyword),
		DictModels:  make(map[int]*EntKeyword),
		Origins:     make(map[Origin]int),
	}
}
