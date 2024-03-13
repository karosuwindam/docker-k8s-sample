package bookmarkfileread

import (
	"book-newread/config"
	"testing"
)

func TestRead(t *testing.T) {
	config.Init()
	if err := Init(); err != nil {
		t.Error(err)
	}
	if urls := ReadBookmark(); len(urls) == 0 {
		t.Error("Not Read Bookmark")
	}
}
