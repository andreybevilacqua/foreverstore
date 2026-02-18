package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpicture"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "6804429f74181a63c50c3d81d733a12f14a353ff"
	expectedPathname := "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"
	assert.Equal(t, expectedPathname, pathKey.PathName)
	assert.Equal(t, expectedOriginalKey, pathKey.Filename)
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momspecials"
	data := []byte("some jpg bytes")
	assert.Nil(t, s.writeStream(key, bytes.NewReader(data)))

	r, err := s.Read(key)
	assert.Nil(t, err)

	b, err := io.ReadAll(r)
	assert.Equal(t, string(data), string(b))
	fmt.Printf("File content: %s", string(b))
	s.Delete(key)
}

func TestStoreDelete(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momspecials"
	data := []byte("some jpg bytes")
	err := s.writeStream(key, bytes.NewReader(data))
	assert.Nil(t, err)
	assert.Nil(t, s.Delete(key))
}

func TestHasKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momspecials"
	data := []byte("some jpg bytes")
	err := s.writeStream(key, bytes.NewReader(data))
	assert.Nil(t, err)
	assert.True(t, s.Has(key))
	assert.Nil(t, s.Delete(key))
}
