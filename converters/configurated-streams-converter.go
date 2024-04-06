package converters

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
)

func ConvertDtoToConfiguratedStreamModify(node dtos.Node, stream *dtos.ConfiguredStream, configuration entities.Configuration) entities.ConfiguredStream {
	return entities.ConfiguredStream{
		ConfiguredStreamId:   stream.ConfiguredStreamId,
		ConfigurationId:      configuration.ConfigurationId,
		NodeId:               node.Id,
		StreamId:             stream.StreamId,
		UpdateFrequency:      stream.UpdateFrequency,
		CheckErrors:          stream.CheckErrors,
		NormalUpperThreshold: stream.UpperThreshold,
		NormalLowerThreshold: stream.LowerThreshold,
		CalibrationId:        stream.CalibrationId,
	}
}

func ConvertDtoToConfiguratedStream(node dtos.Node, stream *dtos.ConfiguredStream, configuration entities.Configuration) entities.ConfiguredStream {
	return entities.ConfiguredStream{
		ConfiguredStreamId:   stream.ConfiguredStreamId,
		ConfigurationId:      configuration.ConfigurationId,
		NodeId:               node.Id,
		StreamId:             stream.StreamId,
		UpdateFrequency:      stream.UpdateFrequency,
		CheckErrors:          stream.CheckErrors,
		NormalUpperThreshold: stream.UpperThreshold,
		NormalLowerThreshold: stream.LowerThreshold,
		CalibrationId:        stream.CalibrationId,
	}
}
