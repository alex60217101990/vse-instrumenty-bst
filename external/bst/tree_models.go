package bst

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sync"
	"sync/atomic"

	"github.com/alex60217101990/vse-instrumenty-bst/external/logger"

	"github.com/alex60217101990/vse-instrumenty-bst/external/configs"
	"github.com/valyala/gozstd"
	"golang.org/x/sync/errgroup"
)

type Node struct {
	key   int
	value interface{}
	left  *Node
	right *Node
}

func (n *Node) Insert(key int, value interface{}) error {

	if n == nil {
		return InsertIntoNilTree
	}

	switch {
	case key == n.key:
		return nil
	case key < n.key:
		if n.left == nil {
			n.left = &Node{key: key, value: value}
			return nil
		}
		return n.left.Insert(key, value)
	case key > n.key:
		if n.right == nil {
			n.right = &Node{key: key, value: value}
			return nil
		}
		return n.right.Insert(key, value)
	}
	return nil
}

func (n *Node) Search(key int) (interface{}, bool) {

	if n == nil {
		return "", false
	}

	switch {
	case key == n.key:
		return n.value, true
	case key < n.key:
		return n.left.Search(key)
	default:
		return n.right.Search(key)
	}
}

func (n *Node) findMax(parent *Node) (*Node, *Node) {
	if n == nil {
		return nil, parent
	}
	if n.right == nil {
		return n, parent
	}
	return n.right.findMax(n)
}

func (n *Node) replaceNode(parent, replacement *Node) error {
	if n == nil {
		return ReplaceNotAllowed
	}

	if n == parent.left {
		parent.left = replacement
		return nil
	}
	parent.right = replacement
	return nil
}

func (n *Node) Delete(key int, parent *Node) error {
	if n == nil {
		return DelNotExistKey
	}

	switch {
	case key < n.key:
		return n.left.Delete(key, n)
	case key > n.key:
		return n.right.Delete(key, n)
	default:
		if n.left == nil && n.right == nil {
			n.replaceNode(parent, nil)
			return nil
		}

		if n.left == nil {
			n.replaceNode(parent, n.right)
			return nil
		}

		if n.right == nil {
			n.replaceNode(parent, n.left)
			return nil
		}

		replacement, replParent := n.left.findMax(n)

		n.key = replacement.key
		n.value = replacement.value

		return replacement.Delete(replacement.key, replParent)
	}
}

// BinarySearchTree the binary search tree of interface
type BinarySearchTree struct {
	root *Node
	sync.RWMutex
	size uint64
}

// Load load tree from file or network
func (bst *BinarySearchTree) Load(r io.ReadCloser, useCompress ...bool) (err error) {
	defer func() {
		if err != nil {
			logger.AppLogger.Errorw(err.Error(), "bst", "load")
		}
	}()

	var decoder *json.Decoder
	bst.root = nil

	if len(useCompress) == 0 {
		decoder = json.NewDecoder(r)

		for decoder.More() {
			err = decoder.Decode(bst)
			if err != nil {
				return err
			}
		}

	} else {
		pr, pw := io.Pipe()
		eg := new(errgroup.Group)

		eg.Go(func() (err error) {
			defer pr.Close()
			decoder = json.NewDecoder(pr)

			for decoder.More() {
				err = decoder.Decode(bst)
				if err != nil {
					return err
				}
			}

			return err
		})

		eg.Go(func() error {
			defer pw.Close()
			return gozstd.StreamDecompress(pw, r)
		})

		err = eg.Wait()
	}

	return err
}

// Dump snapshot tree to file or network
func (bst *BinarySearchTree) Dump(w io.WriteCloser, useCompress ...bool) (err error) {
	defer func() {
		if err != nil {
			logger.AppLogger.Errorw(err.Error(), "bst", "dump")
		}
	}()

	var encoder *json.Encoder

	if len(useCompress) > 0 {
		pr, pw := io.Pipe()
		eg := new(errgroup.Group)

		eg.Go(func() (err error) {
			defer pw.Close()
			encoder = json.NewEncoder(pw)
			encoder.SetIndent("", "\t")
			return encoder.Encode(bst)
		})

		eg.Go(func() error {
			defer pr.Close()
			return gozstd.StreamCompressLevel(w, pr, 22)
		})

		return eg.Wait()
	}

	encoder = json.NewEncoder(w)
	encoder.SetIndent("", "\t")

	return encoder.Encode(bst)
}

// Insert calls `Node.Insert` unless the root node is `nil`
func (bst *BinarySearchTree) Insert(key int, value interface{}, useLock ...bool) (err error) {
	if len(useLock) > 0 {
		bst.Lock()
		defer bst.Unlock()
	}
	defer func() {
		if err == nil {
			atomic.AddUint64(&bst.size, 1)
			if configs.Conf.IsDebug {
				logger.AppLogger.Warnf("current tree size: %d", atomic.LoadUint64(&bst.size))
			}
			return
		}

		logger.AppLogger.Errorw(err.Error(), "bst", "insert")
	}()

	if bst.root == nil {
		bst.root = &Node{key, value, nil, nil}

		return err
	}

	if atomic.LoadUint64(&bst.size) > 30 {
		return TreeBufferLimitExceeded
	}

	return bst.root.Insert(key, value)
}

// Search calls `Node.Find` unless the root node is `nil`
func (bst *BinarySearchTree) Search(key int) (interface{}, bool) {
	if bst.root == nil {
		return "", false
	}
	return bst.root.Search(key)
}

// Delete has one special case: the empty tree. (And deleting from an empty tree is an error.)
// In all other cases, it calls `Node.Delete`.
func (bst *BinarySearchTree) Delete(key int) (err error) {
	defer func() {
		if err == nil {
			atomic.AddUint64(&bst.size, uint64(math.MaxUint64))
			if configs.Conf.IsDebug {
				logger.AppLogger.Warnf("current tree size: %d", atomic.LoadUint64(&bst.size))
			}
			return
		}

		logger.AppLogger.Errorw(err.Error(), "bst", "delete")
	}()

	if bst.root == nil {
		return DelFromEmptyTree
	}

	fakeParent := &Node{right: bst.root}
	err = bst.root.Delete(key, fakeParent)
	if err != nil {
		return err
	}

	if fakeParent.right == nil {
		bst.root = nil
	}

	return nil
}

// MarshalJSON custom marshaler for tree
func (bst *BinarySearchTree) MarshalJSON() ([]byte, error) {
	bstMap := make(map[int]interface{})

	if bst == nil {
		bst = &BinarySearchTree{}
	}

	bst.mapBuilder(bst.root, 0, true, bstMap)

	return json.Marshal(bstMap)
}

// UnmarshalJSON custom unmarshaler for tree
func (bst *BinarySearchTree) UnmarshalJSON(data []byte) (err error) {
	bstMap := make(map[int]interface{})

	if err = json.Unmarshal(data, &bstMap); err != nil {
		return err
	}

	if bst == nil {
		bst = &BinarySearchTree{}
	}

	for key, val := range bstMap {
		bst.Insert(key, val, true)
	}

	return err
}

// internal recursive function to build map data a tree
func (bst *BinarySearchTree) mapBuilder(n *Node, level int, useLock bool, serializeData map[int]interface{}) {
	if n != nil {
		level++
		bst.mapBuilder(n.left, level, useLock, serializeData)
		func() {
			if useLock {
				bst.Lock()
				defer bst.Unlock()
			}
			serializeData[n.key] = n.value
		}()
		bst.mapBuilder(n.right, level, useLock, serializeData)
	}
}

// String prints a visual representation of the tree
func (bst *BinarySearchTree) String() {
	bst.Lock()
	defer bst.Unlock()
	fmt.Println("------------------------------------------------")
	stringify(bst.root, 0)
	fmt.Println("------------------------------------------------")
}

// internal recursive function to print a tree
func stringify(n *Node, level int) {
	if n != nil {
		format := ""
		for i := 0; i < level; i++ {
			format += "       "
		}
		format += "---[ "
		level++
		stringify(n.left, level)
		fmt.Printf(format+"%d => %v ]\n", n.key, n.value)
		stringify(n.right, level)
	}
}
