package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
)

type ImageFaceInfo [][]int

type ImageFaceInfoStore map[string]ImageFaceInfo

func StoreNew(path string) (ImageFaceInfoStore, error) {
	file, err := os.Open(path)
	if err != nil {
		return ImageFaceInfoStore{}, nil
	}

	ret := ImageFaceInfoStore{}
	err = json.NewDecoder(file).Decode(&ret)
	return ret, err
}

func (store *ImageFaceInfoStore) Save(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		return err
	}

	defer file.Close()

	b, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	reader := bytes.NewReader(b)
	_, err = io.Copy(file, reader)

	return err
}

func (store *ImageFaceInfoStore) Add(name string, image ImageFaceInfo) {
	if store == nil {
		*store = make(map[string]ImageFaceInfo)
	}
	(*store)[name] = image
}

func (store *ImageFaceInfoStore) Check(path string) {
	for k := range *store {
		if _, err := os.Stat(path + "/" + k); err != nil {
			log.Println(err)
			delete(*store, k)
		}
	}
}

func (store *ImageFaceInfoStore) Jsonify() string {
	b, err := json.MarshalIndent(*store, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}
