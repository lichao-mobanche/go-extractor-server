package controllers

import (
	
	"net/http"
	
	"github.com/lichao-mobanche/go-extractor-server/server/global"
)

// ErrorCode TODO
func ErrorCode(err error) int {
	var code int
	switch err.(type) {
	case InvalidBody:
		code = http.StatusBadRequest
	case global.QueueUnavailableError:
		code = http.StatusNotAcceptable
	case global.QueueFullError:
		code = http.StatusTooManyRequests
	default:
		code = http.StatusInternalServerError
	}
	return code
}