package controllers

import (
	"fmt"

	"github.com/cfhamlet/os-rq-pod/pkg/sth"
	"github.com/gin-gonic/gin"
	"github.com/lichao-mobanche/go-extractor-server/pkg/request"
	"github.com/lichao-mobanche/go-extractor-server/server/exqueue"
	extract "github.com/lichao-mobanche/go-extractor-server/pkg/extract"
)

// ExQueueController TODO
type ExQueueController struct {
	exq *exqueue.ExQueue
}

// NewExQueueController TODO
func NewExQueueController(e *exqueue.ExQueue) *ExQueueController {
	return &ExQueueController{e}
}

// ExtractLinks TODO
func (ctrl *ExQueueController) ExtractLinks(c *gin.Context) (res sth.Result, err error) {

	var req *request.Request = &request.Request{Exfunc:extract.LinkExtract}
	if err = c.ShouldBindJSON(req); err != nil {
		err = InvalidBody(fmt.Sprintf("%s", err))
		return
	}
	ctrl.exq.AddRequest(req)

	resorerr := <-req.Responsec
	switch resorerr.(type) {
	case error:
		err = resorerr.(error)
	case request.Response:
		res = resorerr.(request.Response)
	}

	return
}

// Extract TODO
func(ctrl *ExQueueController) Extract(c *gin.Context) (res sth.Result, err error) {

	var req *request.ExtractRequest = &request.ExtractRequest{Exfunc:extract.Extract}
	if err = c.ShouldBindJSON(req); err != nil {
		err = InvalidBody(fmt.Sprintf("%s", err))
		return
	}
	ctrl.exq.AddRequest(req)

	resorerr := <-req.Responsec
	switch resorerr.(type) {
	case error:
		err = resorerr.(error)
	case request.Response:
		res = resorerr.(request.Response)
	}

	return
}
