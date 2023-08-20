package goentdb

import (
	"errors"
	"fmt"
	"html"
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

	title := strings.ToLower(video.Title)
	title = strings.Replace(title, "-", " ", -1)
	tokens := strings.Split(title, " ")

	for _, token := range tokens {
		token = strings.Trim(token, TrimSymbols)
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

func (edb *EntDB) GetKeywordsRelatedSet(Video *EntVideo, Size int, UseSeoPool bool, Exclude []*EntVideo) []*EntKeyword {
	res := make([]*EntKeyword, 0)

	videos, _ := edb.RelevantBySearch(Video.Slug, 300)

	if len(videos) < 20 {
		return edb.GetKeywordsRandomSet(Size, UseSeoPool, Exclude)
	}

	cur := 0
	maxCur := 1000
	for len(res) < Size {
		ev := videos[cur%len(videos)]
		// ev := edb.Random()
		// TODO: if (ev not in Exclude and ev not in Taken) {
		// Check if EntVideo has any associated EntKeyword
		if len(ev.Keywords) > 0 {
			res = append(res, ev.GetRandomKeyword())
		}
		cur++
		//}
		if cur == maxCur {
			break
		}
	}

	return res
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

/*
	GetRelevantForVideo
	===================
	relevant videos with similar tags and models
*/
func (edb *EntDB) GetRelevantForVideo(Video *EntVideo, Size int) ([]*EntVideo, int) {
	Counter := make(map[*EntVideo]int)

	for _, tag := range Video.Tags {
		fmt.Printf("tag: %v\n", tag)
		videos, _ := edb.RandomSetByTag(tag.GetSlug(), Size)
		for _, video := range videos {
			Counter[video]++
		}
	}

	for _, model := range Video.Models {
		fmt.Printf("model: %v\n", model)
		videos, _ := edb.RandomSetByModel(model.GetSlug(), Size)
		for _, video := range videos {
			Counter[video]++
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

/*
	Get slice of random EntVideos based on query filter
*/
func (edb *EntDB) RelevantBySearch(Slug string, Size int) ([]*EntVideo, int) {

	Weights := map[string]int{
		"video": 0,
		"porn":  0,
		"sex":   0,
		"fuck":  0,
		"xxx":   0,
		"xnxx":  0,
	}

	QueryTokens := strings.Split(strings.ToLower(Slug), "-")
	Counter := make(map[*EntVideo]int)

	for _, token := range QueryTokens {
		token = strings.Trim(token, TrimSymbols)

		if _, exists := StopWordsMap[token]; exists {
			continue
		}

		if len(token) < 3 {
			continue
		}

		if videos, exists := edb.Search[token]; exists {
			for _, video := range videos {
				if weight, found := Weights[token]; found {
					Counter[video] = Counter[video] + weight
				} else {
					Counter[video]++
				}
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

/*
	Get RelevantVideos for EntVideo
	Exclude MainVideo from the result
*/
func (edb *EntDB) GetRelevantForVideoBySearch(Video *EntVideo, Size int) ([]*EntVideo, int) {
	Title := html.UnescapeString(Video.Title)
	QueryTokens := strings.Split(strings.ToLower(Title), " ")
	Counter := make(map[*EntVideo]int)

	for _, token := range QueryTokens {
		token = strings.Trim(token, TrimSymbols)

		if _, exists := StopWordsMap[token]; exists {
			continue
		}

		if len(token) < 3 {
			continue
		}

		if videos, exists := edb.Search[token]; exists {
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
		if k.Id == Video.Id {
			continue
		}
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

	// if expected size bigger than actual list then just take a list
	if Size >= len(models) {
		return models, len(models)
	}

	seen := make(map[*EntVideo]bool)
	res := make([]*EntVideo, Size)
	taken := 0

	// While taken less than expected and didn't see every item
	for taken < Size && len(seen) < len(models) {
		candidate = models[rand.Intn(len(models))]
		if candidate == nil {
			continue
		}
		if _, exists := seen[candidate]; !exists {
			res[taken] = candidate
			seen[candidate] = true
			taken++
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
	taken := 0
	for taken < Size && len(seen) < len(tags) {
		candidate = tags[rand.Intn(len(tags))]
		if candidate == nil {
			continue
		}
		if _, exists := seen[candidate]; !exists {
			res[taken] = candidate
			seen[candidate] = true
			taken++
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
	res := make([]*EntKeyword, 0)

	for len(res) < Size {
		kw := edb.Random().GetRandomKeyword()
		if kw.Phrase != "not-found" {
			res = append(res, kw)
		}
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
