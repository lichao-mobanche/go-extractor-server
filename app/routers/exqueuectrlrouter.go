package routers

import (
	"github.com/lichao-mobanche/go-extractor-server/app/controllers"
	"github.com/cfhamlet/os-rq-pod/pkg/ginserv"
	"github.com/cfhamlet/os-rq-pod/pkg/ginserv/route"
)

// RouteExQueueCtrl TODO
func RouteExQueueCtrl(root ginserv.RouterGroup, ctrl *controllers.ExQueueController) {
	g := root.Group("/")
	routes := []*route.Route{
		route.New(g.POST, "/", ctrl.ExtractLinks),
	}
	route.Bind(routes, controllers.ErrorCode)
}