package dtos

type MetricCard struct {
	Name  string
	Value float64
}

func NewMetricCard(name string, value float64) MetricCard {
	return MetricCard{Name: name, Value: value}
}
