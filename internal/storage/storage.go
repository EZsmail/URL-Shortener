package storage

import (
	"errors"
)

var (
	ErrURLNotFound   = errors.New("url not found")
	ErrURLExists     = errors.New("url exist")
	ErrAliasNotFound = errors.New("alias not found")
)
