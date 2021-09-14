package main

import (
	"fmt"
	restful "github.com/emicklei/go-restful/v3"
	"net/http"
	"strings"
)

func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/temperature").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	service.Route(service.GET("/all").
		To(AllEstimations))

	service.Route(service.GET("/{name}").
		To(ParticularEstimation))

	return service
}

func AllEstimations(request *restful.Request, response *restful.Response) {
	estimations, err := EstimateAll()
	if err != nil {
		processErrors(response, err)
	} else {
		response.WriteEntity(estimations)
	}
}

func ParticularEstimation(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")
	estimation, err := Estimate(name)
	if err != nil {
		processErrors(response, err)
	} else {
		response.WriteEntity(estimation)
	}
}

func processErrors(response *restful.Response, err error) {
	formatted := fmt.Sprintf("Error occured: %s.", strings.Title(err.Error()))

	switch err.Error() {
		case ErrARSOUnknownLocation:
			response.WriteErrorString(http.StatusNotFound, formatted)
		case ErrARSOUnreachable:
		case ErrARSOBadResponse:
		case ErrARSOResponseNotParsable:
			response.WriteErrorString(http.StatusInternalServerError, formatted)
	}
}
