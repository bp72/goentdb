package goentdb

import (
	"testing"
)

func TestEntKeyword(t *testing.T) {
	var keyword EntKeyword

	keyword.Phrase = "Abc EfD HIJ"

	Expected := "abc-efd-hij"
	Got := keyword.GetSlug()

	if Got != Expected {
		t.Errorf("test slug uri failed: got %v, wanted %v", Got, Expected)
	}

	keyword.Phrase = "#Abc EfD HIJ"

	Expected = "abc-efd-hij"
	Got = keyword.GetSlug()

	if Got != Expected {
		t.Errorf("test slug uri failed: got %v, wanted %v", Got, Expected)
	}
}
