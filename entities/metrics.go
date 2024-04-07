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
		return "Máximo"
	case Minimo:
		return "Mínimo"
	case Nulos:
		return "Cantidad de Nulos"
	case Observaciones:
		return "Observaciones"
	case AguasAlerta:
		return "Cantidad sobre nivel de Alerta"
	case AguasEvacuacion:
		return "Cantidad sobre nivel de Evacuación"
	case AguasBajas:
		return "Cantidad bajo nivel de Aguas Bajas"
	default:
		return "Unknown"
	}
}
