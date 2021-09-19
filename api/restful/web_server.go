package main

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	temperatures "github.com/klokar/kolpa-temperatures"
	"github.com/klokar/kolpa-temperatures/arso"
	"net/http"
	"strings"
)

const (
	BaseURL = "/temperature"
)

type RestfulWebRouterService struct {
	BaseUrl               string
	MimeType              string
	TemperaturesEstimator temperatures.TemperatureEstimator
}

func (router RestfulWebRouterService) Initialize() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path(router.BaseUrl).
		Consumes(router.MimeType).
		Produces(router.MimeType)

	service.Route(service.GET("/all").
		To(router.All))

	service.Route(service.GET("/{name}").
		To(router.Single))

	return service
}

func (router RestfulWebRouterService) All(request *restful.Request, response *restful.Response) {
	estimations, err := router.TemperaturesEstimator.All()
	if err != nil {
		processErrors(response, err)
	} else {
		response.WriteEntity(estimations)
	}
}

func (router RestfulWebRouterService) Single(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")
	estimation, err := router.TemperaturesEstimator.Single(name)
	if err != nil {
		processErrors(response, err)
	} else {
		response.WriteEntity(estimation)
	}
}

func processErrors(response *restful.Response, err error) {
	formatted := fmt.Sprintf("Error occured: %s.", strings.Title(err.Error()))

	switch err.Error() {
	case arso.ErrARSOUnknownLocation:
		response.WriteErrorString(http.StatusNotFound, formatted)
	case arso.ErrARSOUnreachable:
	case arso.ErrARSOBadResponse:
	case arso.ErrARSOResponseNotParsable:
		response.WriteErrorString(http.StatusInternalServerError, formatted)
	}
}
