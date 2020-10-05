package bst

import "io"

// Tree interface for implementation of binary search tree data struct
type Tree interface {
	Load(r io.ReadCloser, useCompress ...bool) error
	Dump(w io.WriteCloser, useCompress ...bool) error
	Insert(key int, value interface{}, useLock ...bool) error
	Search(key int) (interface{}, bool)
	Delete(key int) error
	String()
}
