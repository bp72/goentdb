package goentdb

import (
	"strings"
)

type EntSlug string

func (self EntSlug) GetSlug() string {
	return strings.Replace(strings.Replace(strings.ToLower(string(self)), " ", "-", -1), "#", "", -1)
}
