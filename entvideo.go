package goentdb

import (
	"bytes"
	"fmt"
	"html"
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
	MapKeyword map[string]*EntKeyword // Ad-hoc. TODO review later
	Owner      *EntDB                 // Ad-hoc. TODO review later
}

func (ev *EntVideo) GetTitle() string {
	return html.UnescapeString(ev.Title)
}

func (ev *EntVideo) IsRefererDisabled() bool {
	return ev.Origin == OriginAlphaporno
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
	return fmt.Sprintf("%s/%s", ev.Owner.ThumbBaseUrl, ev.GetPosterThumbRelatedPath())
}

func (ev *EntVideo) GetSlug() string {
	lc := strings.ToLower(ev.Title)
	var buffer bytes.Buffer

	for _, char := range lc {
		if (char-'a' >= 0 && char-'a' < 26) || (char-'0' >= 0 && char-'0' < 10) || (char-' ' == 0) {
			if char == ' ' {
				buffer.WriteRune('-')
			} else {
				buffer.WriteRune(char)
			}
		}
	}

	return buffer.String()
}

func (ev *EntVideo) GetMD5() string {
	return MD5(ev.Slug)
}

func (ev *EntVideo) GetRandomKeyword() *EntKeyword {
	if len(ev.Keywords) > 0 {
		keyword := ev.Keywords[rand.Intn(len(ev.Keywords))]
		return keyword
	}
	return NewKeyword(0, "not-found")
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
	return nil, fmt.Errorf("EntVideo %d does not have %s", ev.Id, Slug)
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

func NewEntVideo(edb *EntDB) *EntVideo {
	return &EntVideo{
		MapKeyword: make(map[string]*EntKeyword),
		Owner:      edb,
	}
}

func (v *EntVideo) GetTitleTokens(excludeStopWords bool) []string {
	lcTitle := strings.ToLower(v.GetTitle())
	var buffer bytes.Buffer

	res := make([]string, 0)

	for _, char := range lcTitle {
		if (char-'a' >= 0 && char-'a' < 26) || (char-'0' >= 0 && char-'0' < 10) || (char-' ' == 0) {
			if char == ' ' {
				res = append(res, buffer.String())
				if excludeStopWords {
					if _, exists := StopWordsMap[res[len(res)-1]]; exists {
						res = res[:len(res)-1]
					}
				}
				buffer.Reset()
			} else {
				buffer.WriteRune(char)
			}
		}
	}
	res = append(res, buffer.String())
	if excludeStopWords {
		if _, exists := StopWordsMap[res[len(res)-1]]; exists {
			res = res[:len(res)-1]
		}
	}
	return res
}

func (v *EntVideo) GetNGrams(useStopWords bool) ([]string, []string) {
	tokens := v.GetTitleTokens(useStopWords)
	tokens2Gram := make([]string, 0)
	tokens3Gram := make([]string, 0)

	for i := 0; i < len(tokens)-1; i++ {
		tokens2Gram = append(tokens2Gram, tokens[i]+" "+tokens[i+1])
	}

	for i := 0; i < len(tokens)-2; i++ {
		tokens3Gram = append(tokens3Gram, tokens[i]+" "+tokens[i+1]+" "+tokens[i+2])
	}

	return tokens2Gram, tokens3Gram
}
