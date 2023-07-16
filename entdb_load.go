package goentdb

import (
	"encoding/gob"
	"os"
)

func (edb *EntDB) LoadTags() error {
	return LoadMapFromFilepath(edb.GetDictTagsPath(), &edb.DictTags, edb.lock)
}

func (edb *EntDB) LoadModels() error {
	return LoadMapFromFilepath(edb.GetDictModelsPath(), &edb.DictModels, edb.lock)
}

func (edb *EntDB) LoadVideos() error {
	f, err := os.Open(edb.GetDictVideosPath())
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(f)

	//edb.lock.RLock()
	//defer edb.lock.RUnlock()

	items := make([]EntVideoForLoad, 0)
	decoder.Decode(&items)

	for _, v := range items {
		edb.AddVideoFromLoad(&v)
	}

	f.Close()
	return nil
}

func (edb *EntDB) Load() {
	edb.LoadTags()
	edb.LoadModels()
	edb.LoadVideos()
}
