package main

import (
	"errors"
)

const (
	ARSOPetrina = 4820
	ARSOMetlika = 4860
	ErrARSOUnknownLocation = "unknown location"
	ErrARSOBaseTempsMissing = "base ARSO temperatures missing"
)

type location struct {
	offset   float64
	fullName string
}

type baseTemperatures struct {
	start float64
	end   float64
}

type Estimation struct {
	FullName    string
	Temperature float64
}

func Estimate(name string) (Estimation, error) {
	var estimation Estimation
	base, err := obtainBase()
	if err != nil {
		return estimation, err
	}

	locations := locations()
	if location, ok := locations[name]; ok {
		estimation.FullName = location.fullName
		estimation.Temperature = location.calculate(base)

		return estimation, nil
	}

	return estimation, errors.New(ErrARSOUnknownLocation)
}

func EstimateAll() ([]Estimation, error) {
	var estimations []Estimation
	base, err := obtainBase()
	if err != nil {
		return estimations, err
	}

	locations := locations()
	for _, location := range locations {
		estimation := Estimation{location.fullName, location.calculate(base)}
		estimations = append(estimations, estimation)
	}

	return estimations, nil
}

func obtainBase() (baseTemperatures, error) {
	var base baseTemperatures
	measured, err := Measurements()

	if err != nil {
		return base, err
	}

	for _, station := range measured {
		if station.Id == ARSOPetrina {
			base.start = station.Temperature
		}
		if station.Id == ARSOMetlika {
			base.end = station.Temperature
		}
	}

	if base.start == 0 || base.end == 0 {
		return base, errors.New(ErrARSOBaseTempsMissing)
	}

	return base, nil
}

func locations() map[string]location {
	return map[string]location{
		"petrina":   {0.0, "Mejni prehod Petrina"},
		"kostel":    {0.2, "Pri Kostelu"},
		"sodevci":   {0.5, "Prelesje / Sodevci"},
		"vinica":    {1.1, "Kopališče Vinica"},
		"adlesici":  {1.56, "Kopališče Adlešiči"},
		"podbrezje": {1.62, "Kopališče Podbrežje"},
		"griblje":   {1.74, "Kopališče Griblje"},
		"krasinec":  {1.8, "Kopališče Krasinec"},
		"podzemelj": {1.8, "Kamp Podzemelj"},
		"primostek": {1.1, "Kamp Primostek"},
		"krizevska": {1.1, "Kamp Križevska vas"},
		"metlika":   {1.0, "Mejni prehod Metlika"},
	}
}

func (location location) calculate(base baseTemperatures) float64 {
	difference := base.end - base.start

	return base.start + difference*location.offset
}
