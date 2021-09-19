package kolpa_temperatures

type Temperature struct {
	Location  string
	Estimated bool
	Value     float64
}

type TemperatureEstimator interface {
	Single(name string) (Temperature, error)
	All() ([]Temperature, error)
}
