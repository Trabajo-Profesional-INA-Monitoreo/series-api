package entities

type StreamType uint64

const (
	Observed StreamType = iota
	Forecasted
	Curated
)
