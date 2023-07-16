package goentdb

import (
	"strings"
)

type EntKeyword struct {
	Type   EntKeywordType
	Id     int
	Phrase string
	slug   string
}

func (ek *EntKeyword) GetMD5() string {
	return MD5(ek.GetSlug())
}

func (ek *EntKeyword) GetSlug() string {
	if ek.slug == "" {
		ek.slug = strings.Replace(strings.Replace(strings.ToLower(ek.Phrase), " ", "-", -1), "#", "", -1)
	}
	return ek.slug
}

func NewEntKeyword(id int, phrase string, enttype EntKeywordType) *EntKeyword {
	return &EntKeyword{
		Type:   enttype,
		Id:     id,
		Phrase: phrase,
	}
}

func NewTag(id int, phrase string) *EntKeyword {
	return NewEntKeyword(id, phrase, EntKeywordTag)
}

func NewModel(id int, phrase string) *EntKeyword {
	return NewEntKeyword(id, phrase, EntKeywordModel)
}

func NewKeyword(id int, phrase string) *EntKeyword {
	return NewEntKeyword(id, phrase, EntKeywordKeyword)
}
