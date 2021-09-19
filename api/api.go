package api

import (
	"github.com/emicklei/go-restful/v3"
)

// WebRouterService TODO: this should not contain restful dependency
type WebRouterService interface {
	Initialize() *restful.WebService
	Single(request *restful.Request, response *restful.Response)
	All(request *restful.Request, response *restful.Response)
}
