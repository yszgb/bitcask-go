package bitcaskgo

import "errors"

var (
	ErrKeyIsEmpty = errors.New("key is empty")
	ErrIndexUpdateFailed = errors.New("failed to update index")
)
