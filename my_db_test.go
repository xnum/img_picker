package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNullNew(t *testing.T) {
	_, err := StoreNew("/tmp/no_exists")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSave(t *testing.T) {
	const path = "/tmp/test_save"
	os.Remove(path)
	store, err := StoreNew(path)
	if err != nil {
		t.Fatal(err)
	}

	store.Add("a", ImageFaceInfo{})
	assert.Equal(t, len(store), 1, "add failed")

	store.Add("b", ImageFaceInfo{{1, 2, 3, 4}})
	assert.Equal(t, len(store), 2, "add failed")

	err = store.Save(path)
	if err != nil {
		t.Fatal(err)
	}

	st2, err := StoreNew(path)
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, store, st2, "object not same")

	st2.Check(path)
	assert.Equal(t, len(st2), 0, "st2 has any object!!!")
}
