package app

import "errors"

var (
	ErrorRequest = errors.New("request error")
	ErrorServer  = errors.New("server error")
)
