package entities

type ErrorType uint64

const (
	NullValue ErrorType = iota
	Missing4DaysHorizon
	OutsideOfErrorBands
	ForecastMissing
)
