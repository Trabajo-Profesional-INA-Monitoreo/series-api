package dtos

type MetricCard struct {
	Name  string  `json:"Name"`
	Value float64 `json:"Value"`
}

func NewMetricCard(name string, value float64) MetricCard {
	return MetricCard{Name: name, Value: value}
}
