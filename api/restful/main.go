package main

import (
	"github.com/emicklei/go-restful/v3"
	temperature "github.com/klokar/kolpa-temperatures"
	"github.com/klokar/kolpa-temperatures/api"
	"github.com/klokar/kolpa-temperatures/arso"
	"log"
	"net/http"
)

func main() {
	var measurer arso.TemperatureMeasurer
	measurer = arso.WaterReportMeasurer{ReportURL: arso.WaterReportURL}

	var estimator temperature.TemperatureEstimator
	estimator = arso.TwoBaseEstimator{StartLocationId: arso.LocationIdPetrina, EndLocationId: arso.LocationIdMetlika, Measurer: measurer}

	var webService api.WebRouterService
	webService = RestfulWebRouterService{BaseUrl: BaseURL, MimeType: restful.MIME_JSON, TemperaturesEstimator: estimator}

	restful.Add(webService.Initialize())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
