package entities

type ConfiguredMetric struct {
	MetricId           Metric
	ConfiguredStreamId uint64
}

type Metric uint64

const (
	Mediana Metric = iota // Configurable
	Media                 // Configurable
	Maximo                // Configurable
	Minimo                // Configurable
	Nulos                 // Configurable
	Observaciones
	AguasAlerta
	AguasEvacuacion
	AguasBajas
)

func MapMetricToString(m Metric) string {
	switch m {
	case Mediana:
		return "Mediana"
	case Media:
		return "Media"
	case Maximo:
		return "Maximo"
	case Minimo:
		return "Minimo"
	case Nulos:
		return "Nulos"
	case Observaciones:
		return "Observaciones"
	case AguasAlerta:
		return "AguasAlerta"
	case AguasEvacuacion:
		return "AguasEvacuacion"
	case AguasBajas:
		return "AguasBajas"
	default:
		return "Unknown"
	}
}
