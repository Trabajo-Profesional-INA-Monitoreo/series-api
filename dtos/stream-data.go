package dtos

import "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"

type StreamData struct {
	Metrics         *[]MetricCard
	Station         string
	VarId           uint64
	VarName         string
	AlertLevel      *float64
	EvacuationLevel *float64
	LowWaterLevel   *float64
	Owner           string
	Unit            string
	UnitId          uint64
	Procedure       string
	ProcId          uint64
	UpdateFrequency float64
	StreamType      uint64
}

func NewStreamData(stream entities.Stream, configured entities.ConfiguredStream) *StreamData {
	return &StreamData{
		Station:         stream.Station.Name,
		VarId:           stream.VariableId,
		VarName:         stream.Variable.Name,
		StreamType:      stream.ProcedureId,
		AlertLevel:      &stream.Station.AlertLevel,
		EvacuationLevel: &stream.Station.EvacuationLevel,
		LowWaterLevel:   &stream.Station.LowWaterLevel,
		Unit:            stream.Unit.Name,
		UnitId:          stream.UnitId,
		Procedure:       stream.Procedure.Name,
		ProcId:          stream.ProcedureId,
		Owner:           stream.Station.Owner,
		UpdateFrequency: configured.UpdateFrequency,
	}
}
