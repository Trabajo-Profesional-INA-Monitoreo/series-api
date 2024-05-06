package dtos

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"time"
)

type StreamData struct {
	Metrics                 *[]MetricCard       `json:"Metrics"`
	Station                 string              `json:"Station"`
	VarId                   uint64              `json:"VarId"`
	VarName                 string              `json:"VarName"`
	AlertLevel              *float64            `json:"AlertLevel"`
	EvacuationLevel         *float64            `json:"EvacuationLevel"`
	LowWaterLevel           *float64            `json:"LowWaterLevel"`
	Owner                   string              `json:"Owner"`
	Unit                    string              `json:"Unit"`
	UnitId                  uint64              `json:"UnitId"`
	Procedure               string              `json:"Procedure"`
	ProcId                  uint64              `json:"ProcId"`
	UpdateFrequency         float64             `json:"UpdateFrequency"`
	StreamType              entities.StreamType `json:"StreamType"`
	CalibrationId           uint64              `json:"CalibrationId"`
	ObservedRelatedStreamId *uint64             `json:"ObservedRelatedStreamId"`
	LastUpdate              *time.Time          `json:"LastUpdate"`
	NormalUpperThreshold    float64             `json:"NormalUpperThreshold"`
	NormalLowerThreshold    float64             `json:"NormalLowerThreshold"`
}

func NewStreamData(stream entities.Stream, configured entities.ConfiguredStream) *StreamData {
	return &StreamData{
		Station:                 stream.Station.Name,
		VarId:                   stream.VariableId,
		VarName:                 stream.Variable.Name,
		StreamType:              stream.StreamType,
		AlertLevel:              stream.Station.AlertLevel,
		EvacuationLevel:         stream.Station.EvacuationLevel,
		LowWaterLevel:           stream.Station.LowWaterLevel,
		Unit:                    stream.Unit.Name,
		UnitId:                  stream.UnitId,
		Procedure:               stream.Procedure.Name,
		ProcId:                  stream.ProcedureId,
		Owner:                   stream.Station.Owner,
		UpdateFrequency:         configured.UpdateFrequency,
		CalibrationId:           configured.CalibrationId,
		ObservedRelatedStreamId: configured.ObservedRelatedStreamId,
		NormalUpperThreshold:    configured.NormalUpperThreshold,
		NormalLowerThreshold:    configured.NormalLowerThreshold,
	}
}
