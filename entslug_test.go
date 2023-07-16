package goentdb

import (
	"testing"
)

func TestEntSlug(t *testing.T) {
	var slug EntSlug

	slug = "Abc EfD HIJ"

	Expected := "abc-efd-hij"
	Got := slug.GetSlug()

	if Got != Expected {
		t.Errorf("test slug uri failed: got %v, wanted %v", Got, Expected)
	}

	slug = "#Abc EfD HIJ"

	Expected = "abc-efd-hij"
	Got = slug.GetSlug()

	if Got != Expected {
		t.Errorf("test slug uri failed: got %v, wanted %v", Got, Expected)
	}
}
