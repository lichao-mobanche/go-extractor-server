package routers

import (
	"github.com/cfhamlet/os-rq-pod/pkg/ginserv"
	"github.com/cfhamlet/os-rq-pod/pkg/ginserv/route"
	"github.com/lichao-mobanche/go-extractor-server/app/controllers"
)

// RouteExQueueCtrl TODO
func RouteExQueueCtrl(root ginserv.RouterGroup, ctrl *controllers.ExQueueController) {
	g := root.Group("/")
	routes := []*route.Route{
		route.New(g.POST, "/", ctrl.ExtractLinks),
		route.New(g.POST, "/", ctrl.Extract),
	}
	route.Bind(routes, controllers.ErrorCode)
}
