package storage

import "errors"

var (
	ErrorAliasNotFound = errors.New("alias not found")
	ErrorAliasExists   = errors.New("alias exists")
)
