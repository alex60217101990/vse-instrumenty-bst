package http_server

import "errors"

var (
	EmptyKey      = errors.New("empty 'key' query parameter")
	ValueNotFound = errors.New("value by key not found")
)
