package main

import (
	"encoding/xml"
	"errors"
	memoryCache "github.com/patrickmn/go-cache"
	"io"
	"net/http"
	"time"
)

const (
	CacheExpirationMinutes     = 30
	CacheKey                   = "arso-report"
	UrlARSO                    = "http://www.arso.gov.si/xml/vode/hidro_podatki_zadnji.xml"
	ErrARSOUnreachable         = "ARSO is unreachable"
	ErrARSOBadResponse         = "ARSO returned bad reponse"
	ErrARSOResponseNotParsable = "ARSO response could not be parsed"
)

var cache *memoryCache.Cache

func init() {
	cache = memoryCache.New(CacheExpirationMinutes*time.Minute, CacheExpirationMinutes*2*time.Minute)
}

type Station struct {
	Id          int     `xml:"sifra,attr"`
	River       string  `xml:"reka"`
	Location    string  `xml:"merilno_mesto"`
	Temperature float64 `xml:"temp_vode"`
}

type arsoHydroResponse struct {
	XMLName  xml.Name  `xml:"arsopodatki"`
	Date     string    `xml:"datum_priprave"`
	Stations []Station `xml:"postaja"`
}

func Measurements() ([]Station, error) {
	var measurements []Station
	var results arsoHydroResponse

	report, found := cache.Get(CacheKey)
	if found {
		return report.([]Station), nil
	}

	resp, err := http.Get(UrlARSO)
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

	cache.Set(CacheKey, results.Stations, memoryCache.DefaultExpiration)

	return results.Stations, nil
}
