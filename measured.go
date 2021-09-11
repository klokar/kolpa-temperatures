package main

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
)

const ARSOUrl = "http://www.arso.gov.si/xml/vode/hidro_podatki_zadnji.xml"

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

	resp, err := http.Get(ARSOUrl)
	if err != nil {
		return measurements, errors.New("ARSO unreachable")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return measurements, errors.New("ARSO response unreadable")
	}

	err = xml.Unmarshal(body, &results)
	if err != nil {
		return measurements, errors.New("ARSO response unparsable")
	}

	return results.Stations, nil
}
