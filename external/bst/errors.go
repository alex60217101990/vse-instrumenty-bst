package bst

import "errors"

var (
	DelFromEmptyTree = errors.New("Cannot delete from an empty tree")
	DelNotExistKey   = errors.New("Value to be deleted does not exist in the tree")

	InsertIntoNilTree       = errors.New("Cannot insert a value into a nil tree")
	TreeBufferLimitExceeded = errors.New("Tree limit buffer size exeeded")

	ReplaceNotAllowed = errors.New("replace node not allowed on a nil node")
)
