package arso

import (
	"encoding/xml"
	"errors"
	memoryCache "github.com/patrickmn/go-cache"
	"io"
	"net/http"
	"time"
)

const (
	WaterReportURL             = "http://www.arso.gov.si/xml/vode/hidro_podatki_zadnji.xml"
	CacheExpirationMinutes     = 30
	CacheKey                   = "arso-report"
	ErrARSOUnreachable         = "ARSO is unreachable"
	ErrARSOBadResponse         = "ARSO returned bad reponse"
	ErrARSOResponseNotParsable = "ARSO response could not be parsed"
)

var cache *memoryCache.Cache

func init() {
	cache = memoryCache.New(CacheExpirationMinutes*time.Minute, CacheExpirationMinutes*2*time.Minute)
}

type arsoHydroStation struct {
	Id          int     `xml:"sifra,attr"`
	River       string  `xml:"reka"`
	Location    string  `xml:"merilno_mesto"`
	Temperature float64 `xml:"temp_vode"`
}

type arsoHydroResponse struct {
	XMLName  xml.Name           `xml:"arsopodatki"`
	Date     string             `xml:"datum_priprave"`
	Stations []arsoHydroStation `xml:"postaja"`
}

type WaterReportMeasurer struct {
	ReportURL string
}

func (measurer WaterReportMeasurer) Rivers() ([]RiverMeasurement, error) {
	var measurements []RiverMeasurement
	var results arsoHydroResponse

	report, found := cache.Get(CacheKey)
	if found {
		return report.([]RiverMeasurement), nil
	}

	resp, err := http.Get(measurer.ReportURL)
	if err != nil {
		return measurements, errors.New(ErrARSOUnreachable)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return measurements, errors.New(ErrARSOBadResponse)
	}

	err = xml.Unmarshal(body, &results)
	if err != nil {
		return measurements, errors.New(ErrARSOResponseNotParsable)
	}

	// Transform to river measurements.
	for _, station := range results.Stations {
		measurements = append(measurements, toRiverMeasurement(station))
	}

	cache.Set(CacheKey, measurements, memoryCache.DefaultExpiration)

	return measurements, nil
}

func toRiverMeasurement(station arsoHydroStation) RiverMeasurement {
	return RiverMeasurement{
		Id:          station.Id,
		River:       station.River,
		Location:    station.Location,
		Temperature: station.Temperature,
	}
}
