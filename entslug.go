package goentdb

import (
	"strings"
)

type EntSlug string

func (s EntSlug) GetSlug() string {
	return strings.Replace(strings.Replace(strings.ToLower(string(s)), " ", "-", -1), "#", "", -1)
}
