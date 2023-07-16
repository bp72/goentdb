package goentdb

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type EntVideoForLoad struct {
	Id         uint
	Title      string
	Origin     Origin
	OriginId   string
	OriginUrl  string
	Duration   int
	Slug       string
	Source     string
	Descr      string
	ModifiedAt time.Time
	Tags       []int
	Models     []int
	Keywords   []*EntKeyword
	ThumbUrls  []string
	VideoUrls  []string
}

type EntVideo struct {
	Id         uint
	Title      string
	Origin     Origin
	OriginId   string
	OriginUrl  string
	Duration   int
	Slug       string
	Source     string
	Descr      string
	ModifiedAt time.Time
	Tags       []*EntKeyword
	Models     []*EntKeyword
	Keywords   []*EntKeyword
	ThumbUrls  []string
	VideoUrls  []string
	MapKeyword map[string]*EntKeyword
}

func (ev *EntVideo) IsRefererDisabled() bool {
	if ev.Origin == OriginAlphaporno {
		return true
	}
	return false
}

func (ev *EntVideo) IsEmbed() bool {
	return (ev.Origin == OriginEporner || ev.Origin == OriginPorntube || ev.Origin == OriginSuperporn)
}

func (ev *EntVideo) IsStream() bool {
	return ev.Origin == OriginXvideos || ev.Origin == OriginPornone
}

func (ev *EntVideo) GetDescr() string {
	return ev.Descr
}

func (ev *EntVideo) GetMetaKeywords() string {
	var ans []string

	for _, tag := range ev.Tags {
		ans = append(ans, tag.Phrase)
	}
	for _, model := range ev.Models {
		ans = append(ans, model.Phrase)
	}

	return strings.Join(ans, ",")
}

func (ev *EntVideo) GetSubdirs() string {
	name := fmt.Sprintf("%d", ev.Id)
	return fmt.Sprintf("%s/%s", name[0:2], name[2:4])
}

func (ev *EntVideo) GetPosterThumbRelatedPath() string {
	switch ev.Origin {
	case OriginXvideos:
		return fmt.Sprintf("%s/%d_0.webp", ev.GetSubdirs(), ev.Id)
	default:
		return fmt.Sprintf("%s/%d_0.webp", ev.GetSubdirs(), ev.Id)
	}
}

func (ev *EntVideo) GetPosterThumb() string {
	return ev.GetPosterThumbRelatedPath()
}

func (ev *EntVideo) GetThumb() string {
	// TODO: depricate this function
	// Backwards compatability
	return fmt.Sprintf("https://localhost:18443/thumbs/%s", ev.GetPosterThumbRelatedPath())
}

func (ev *EntVideo) GetMD5() string {
	return MD5(ev.Slug)
}

func (ev *EntVideo) GetRandomKeyword() *EntKeyword {
	keyword := ev.Keywords[rand.Intn(len(ev.Keywords))]
	return keyword
}

func (ev *EntVideo) AddTag(tag *EntKeyword) {
	ev.Tags = append(ev.Tags, tag)
}

func (ev *EntVideo) AddModel(model *EntKeyword) {
	ev.Models = append(ev.Models, model)
}

func (ev *EntVideo) AddKeyword(kw *EntKeyword) {
	ev.Keywords = append(ev.Keywords, kw)
	ev.MapKeyword[kw.GetSlug()] = kw
}

func (ev *EntVideo) GetDurationHuman() string {
	duration := ev.Duration
	seconds := duration % 60
	duration = duration / 60
	minutes := duration % 60
	hours := duration / 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func (ev *EntVideo) GetKeywordBySlug(Slug string) (*EntKeyword, error) {
	// TODO: make it O(1) instead of O(N) where N is number of Keywords.
	for _, kw := range ev.Keywords {
		if kw.GetSlug() == Slug {
			return kw, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("EntVideo %d does not have %s", ev.Id, Slug))
}

func (ev *EntVideo) ToLoad() EntVideoForLoad {
	evfl := EntVideoForLoad{}
	evfl.Id = ev.Id
	evfl.Title = ev.Title
	evfl.Origin = ev.Origin
	evfl.OriginId = ev.OriginId
	evfl.OriginUrl = ev.OriginUrl
	evfl.Duration = ev.Duration
	evfl.Slug = ev.Slug
	evfl.Source = ev.Source
	evfl.Descr = ev.Descr
	evfl.ModifiedAt = ev.ModifiedAt
	evfl.Keywords = ev.Keywords
	evfl.ThumbUrls = ev.ThumbUrls
	evfl.VideoUrls = ev.VideoUrls

	for _, tag := range ev.Tags {
		evfl.Tags = append(evfl.Tags, tag.Id)
	}
	for _, model := range ev.Models {
		evfl.Models = append(evfl.Models, model.Id)
	}

	return evfl
}

func NewEntVideo() *EntVideo {
	return &EntVideo{
		MapKeyword: make(map[string]*EntKeyword),
	}
}
