package goentdb

import (
	"encoding/gob"
	"fmt"
	"os"
)

/*
	Dump DictTags to file of map[int]*EntKeyword
	TODO: EntKeyword suppose to have prop Videos []*EntVideo
	TODO: Prop Videos should be filled during the LoadVideos proceduce
	TODO: Prop Videos should be excludeded from the DumpTags

*/
func (edb *EntDB) DumpTags() error {
	return DumpMapToFilepath(edb.GetDictTagsPath(), edb.DictTags, edb.lock)
}

func (edb *EntDB) DumpModels() error {
	return DumpMapToFilepath(edb.GetDictModelsPath(), edb.DictModels, edb.lock)
}

func (edb *EntDB) DumpVideos() error {
	f, err := os.Create(edb.GetDictVideosPath())
	if err != nil {
		return err
	}

	items := make([]EntVideoForLoad, len(edb.Items))
	for pos, video := range edb.Items {
		items[pos] = video.ToLoad()
	}

	encoder := gob.NewEncoder(f)

	edb.lock.RLock()
	defer edb.lock.RUnlock()

	if err := encoder.Encode(items); err != nil {
		return err
	}
	f.Close()
	return nil
}

func (edb *EntDB) Dump() error {
	fmt.Printf("dumping Tags=%d\n", len(edb.DictTags))
	edb.DumpTags()
	fmt.Printf("dumping Models=%d\n", len(edb.DictModels))
	edb.DumpModels()
	fmt.Printf("dumping Videos=%d\n", len(edb.Items))
	edb.DumpVideos()
	return nil
}
