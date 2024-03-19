package dtos

import "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"

type StreamData struct {
	Metrics         *[]MetricCard
	Station         string
	TotalCount      *string
	VarName         string
	AlertLevel      *float64
	EvacuationLevel *float64
	LowLevel        *float64
	Network         string
	Unit            string
	Procedure       string
	UpdateFrequency float64
	StreamType      uint64
}

func NewStreamData(stream entities.Stream, configured entities.ConfiguredStream) *StreamData {
	return &StreamData{
		Station:    stream.Station.Name,
		Network:    stream.Network.Name,
		VarName:    stream.Variable,
		StreamType: stream.ProcId,
		//TotalCount:      *string
		//AlertLevel      *float64
		//EvacuationLevel *float64
		//LowLevel        *float64
		//Unit            string
		//Procedure       string
		UpdateFrequency: configured.UpdateFrequency,
	}
}
