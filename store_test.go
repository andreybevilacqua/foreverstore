package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpicture"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "6804429f74181a63c50c3d81d733a12f14a353ff"
	expectedPathname := "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"
	if pathKey.PathName != expectedPathname {
		t.Errorf("have %s want %s", pathKey.PathName, expectedPathname)
	}

	if pathKey.Filename != expectedOriginalKey {
		t.Errorf("have %s want %s", pathKey.Filename, expectedOriginalKey)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momspecials"
	data := []byte("some jpg bytes")
	err := s.writeStream(key, bytes.NewReader(data))
	assert.Nil(t, err)

	r, err := s.Read(key)
	assert.Nil(t, err)

	b, err := io.ReadAll(r)
	assert.Equal(t, string(data), string(b))
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

	err = s.Delete(key)
	assert.Nil(t, err)
}
