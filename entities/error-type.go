package entities

type ErrorType uint64

const (
	NullValue ErrorType = iota
	Missing4DaysHorizon
	OutsideOfErrorBands
	ForecastMissing
	ObservedOutlier
	ForecastOutOfBounds
	Delay
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
	case Delay:
		return "Delay"
	default:
		return "UnknownError"
	}
}

func (t ErrorType) Translate() string {
	switch t {
	case NullValue:
		return "valores nulos"
	case Missing4DaysHorizon:
		return "falta de horizonte a 4 dias"
	case OutsideOfErrorBands:
		return "valores fuera de banda de errores"
	case ForecastMissing:
		return "falta de pronóstico"
	case ObservedOutlier:
		return "outliers observados"
	case ForecastOutOfBounds:
		return "pronóstico fuera de umbrales"
	case Delay:
		return "retardo"
	default:
		return "error desconocido"
	}
}
