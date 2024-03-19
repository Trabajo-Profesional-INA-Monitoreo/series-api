package entities

type ConfiguredMetric struct {
	MetricId           Metric
	ConfiguredStreamId uint64
}

type Metric uint64

const (
	Mediana Metric = iota
	Media
	Maximo
	Minimo
	Nulos
	TasaDeCambio
)
