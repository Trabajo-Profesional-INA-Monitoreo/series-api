package dtos

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"time"
)

type StreamCardsParameters struct {
	StreamId        *uint64
	ConfigurationId uint64
	TimeStart       time.Time
	TimeEnd         time.Time
	VarId           *uint64
	ProcId          *uint64
	StationId       *uint64
	StreamType      *entities.StreamType
	Page            int
	PageSize        int
}
