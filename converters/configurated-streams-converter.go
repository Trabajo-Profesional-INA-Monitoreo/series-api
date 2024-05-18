package converters

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
)

func ConvertDtoToConfiguratedStreamModify(node dtos.Node, stream *dtos.ConfiguredStream, configuration entities.Configuration) entities.ConfiguredStream {
	return entities.ConfiguredStream{
		ConfiguredStreamId:      stream.ConfiguredStreamId,
		ConfigurationId:         configuration.ConfigurationId,
		NodeId:                  node.Id,
		StreamId:                stream.StreamId,
		UpdateFrequency:         stream.UpdateFrequency,
		CheckErrors:             stream.CheckErrors,
		NormalUpperThreshold:    stream.UpperThreshold,
		NormalLowerThreshold:    stream.LowerThreshold,
		CalibrationId:           stream.CalibrationId,
		ForecastedRangeHours:    stream.ForecastedRangeHours,
		ObservedRelatedStreamId: stream.ObservedRelatedStreamId,
	}
}

func ConvertDtoToConfiguratedStream(nodeId uint64, stream *dtos.ConfiguredStreamCreate, configuration entities.Configuration) entities.ConfiguredStream {
	return entities.ConfiguredStream{
		ConfigurationId:         configuration.ConfigurationId,
		NodeId:                  nodeId,
		StreamId:                stream.StreamId,
		UpdateFrequency:         stream.UpdateFrequency,
		CheckErrors:             stream.CheckErrors,
		NormalUpperThreshold:    stream.UpperThreshold,
		NormalLowerThreshold:    stream.LowerThreshold,
		CalibrationId:           stream.CalibrationId,
		ForecastedRangeHours:    stream.ForecastedRangeHours,
		ObservedRelatedStreamId: stream.ObservedRelatedStreamId,
	}
}
