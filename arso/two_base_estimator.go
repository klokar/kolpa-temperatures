package arso

import (
	"errors"
	temperature "github.com/klokar/kolpa-temperatures"
)

const (
	LocationIdPetrina       = 4820
	LocationIdMetlika       = 4860
	ErrARSOUnknownLocation  = "unknown location"
	ErrARSOBaseTempsMissing = "base ARSO temperatures missing"
)

type TwoBaseEstimator struct {
	StartLocationId int
	EndLocationId   int
	Measurer        TemperatureMeasurer
}

type TemperatureMeasurer interface {
	Rivers() ([]RiverMeasurement, error)
}

type RiverMeasurement struct {
	Id          int
	River       string
	Location    string
	Temperature float64
}

type baseTemperatures struct {
	start float64
	end   float64
}

type location struct {
	offset   float64
	fullName string
}

func (estimator TwoBaseEstimator) Single(name string) (temperature.Temperature, error) {
	var temp temperature.Temperature
	base, err := obtainBase(estimator)
	if err != nil {
		return temp, err
	}

	locations := locationOffsets()
	if location, ok := locations[name]; ok {
		temp.Location = location.fullName
		temp.Value = location.estimate(base)

		return temp, nil
	}

	return temp, errors.New(ErrARSOUnknownLocation)
}

func (estimator TwoBaseEstimator) All() ([]temperature.Temperature, error) {
	var temps []temperature.Temperature
	base, err := obtainBase(estimator)
	if err != nil {
		return temps, err
	}

	locations := locationOffsets()
	for _, location := range locations {
		temp := temperature.Temperature{Location: location.fullName, Value: location.estimate(base)}
		temps = append(temps, temp)
	}

	return temps, nil
}

func obtainBase(estimator TwoBaseEstimator) (baseTemperatures, error) {
	var base baseTemperatures
	measured, err := estimator.Measurer.Rivers()

	if err != nil {
		return base, err
	}

	for _, station := range measured {
		if station.Id == estimator.StartLocationId {
			base.start = station.Temperature
		}
		if station.Id == estimator.EndLocationId {
			base.end = station.Temperature
		}
	}

	if base.start == 0 || base.end == 0 {
		return base, errors.New(ErrARSOBaseTempsMissing)
	}

	return base, nil
}

func locationOffsets() map[string]location {
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

func (location location) estimate(base baseTemperatures) float64 {
	difference := base.end - base.start

	return base.start + difference*location.offset
}
