package entities

type ErrorType uint64

const (
	NullValue ErrorType = iota
	Missing4DaysHorizon
	OutsideOfErrorBands
	ForecastMissing
	ObservedOutlier
	ForecastOutOfBounds
)

func MapErrorTypeToString(e ErrorType) string {
	switch e {
	case NullValue:
		return "NullValue"
	case Missing4DaysHorizon:
		return "Missing4DaysHorizon"
	case OutsideOfErrorBands:
		return "OutsideOfErrorBands"
	case ForecastMissing:
		return "ForecastMissing"
	case ObservedOutlier:
		return "ObservedOutlier"
	case ForecastOutOfBounds:
		return "ForecastOutOfBounds"
	default:
		return "UnknownError"
	}
}
