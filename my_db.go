/*
   persistent data type
*/
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

// StoreNew reads data from path and returns the InfoStore
func StoreNew(path string) (ImageFaceInfoStore, error) {
	file, err := os.Open(path)
	if err != nil {
		return ImageFaceInfoStore{}, nil
	}

	ret := ImageFaceInfoStore{}
	err = json.NewDecoder(file).Decode(&ret)
	return ret, err
}

// Save serializes the InfoStore and write to file with given path
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

// Add sets ImageFaceInfo with its name usually is filename
func (store *ImageFaceInfoStore) Add(name string, image ImageFaceInfo) {
	if store == nil {
		*store = make(map[string]ImageFaceInfo)
	}
	(*store)[name] = image
}

// Check removes the stored data that its file is invalid
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
