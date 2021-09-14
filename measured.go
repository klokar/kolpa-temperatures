package main

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
)

const (
	UrlARSO = "http://www.arso.gov.si/xml/vode/hidro_podatki_zadnji.xml"
	ErrARSOUnreachable = "ARSO is unreachable"
	ErrARSOBadResponse = "ARSO returned bad reponse"
	ErrARSOResponseNotParsable = "ARSO response could not be parsed"
)

type Station struct {
	Id          int  `xml:"sifra,attr"`
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

	return results.Stations, nil
}
