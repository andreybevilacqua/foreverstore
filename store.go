package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "ggnetwork"

type PathKey struct {
	PathName string
	Filename string
}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (p PathKey) FullPath() string {
	return concatRootToPath(p.PathName, p.Filename)
}

type StoreOpts struct {
	FolderRoot        string
	PathTransformFunc PathTransformFunc
}

type PathTransformFunc func(string) PathKey

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.FolderRoot) == 0 {
		opts.FolderRoot = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		Filename: key,
	}
}

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := range sliceLen {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		PathName: strings.Join(paths, "/"),
		Filename: hashStr,
	}
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := concatRootToPath(s.FolderRoot, pathKey.FullPath())
	_, err := os.Stat(fullPathWithRoot)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()
	firstPathNameWithRoot := concatRootToPath(s.FolderRoot, pathKey.FirstPathName())
	return os.RemoveAll(firstPathNameWithRoot)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := concatRootToPath(s.FolderRoot, pathKey.FullPath())
	return os.Open(fullPathWithRoot)
}

func (s *Store) writeStream(key string, content io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	directoryTreeWithRoot := concatRootToPath(s.FolderRoot, pathKey.PathName)
	if err := os.MkdirAll(directoryTreeWithRoot, os.ModePerm); err != nil {
		return err
	}
	fullPath := concatRootToPath(s.FolderRoot, pathKey.FullPath())
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, content)
	if err != nil {
		return err
	}
	log.Printf("written (%d) bytes to disk: %s", n, fullPath)
	return nil
}

func concatRootToPath(folderRoot string, path string) string {
	return fmt.Sprintf("%s/%s", folderRoot, path)
}
