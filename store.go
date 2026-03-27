package main

import (
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
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.FolderRoot)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()
	firstPathNameWithRoot := concatRootToPath(s.FolderRoot, pathKey.FirstPathName())
	return os.RemoveAll(firstPathNameWithRoot)
}

func (s *Store) Read(key string) (int64, io.Reader, error) {
	return s.readStream(key)
}

func (s *Store) Write(key string, content io.Reader) (int64, error) {
	return s.writeStream(key, content)
}

func (s *Store) readStream(key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := concatRootToPath(s.FolderRoot, pathKey.FullPath())

	file, err := os.Open(fullPathWithRoot)
	if err != nil {
		return 0, nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}
	return fi.Size(), file, nil
}

func (s *Store) writeStream(key string, content io.Reader) (int64, error) {
	pathKey := s.PathTransformFunc(key)
	directoryTreeWithRoot := concatRootToPath(s.FolderRoot, pathKey.PathName)
	if err := os.MkdirAll(directoryTreeWithRoot, os.ModePerm); err != nil {
		return 0, err
	}
	fullPath := concatRootToPath(s.FolderRoot, pathKey.FullPath())
	f, err := os.Create(fullPath)
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(f, content)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func concatRootToPath(folderRoot string, path string) string {
	return fmt.Sprintf("%s/%s", folderRoot, path)
}
