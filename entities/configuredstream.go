package entities

type ConfiguredStream struct {
	ConfiguredStreamId uint64 `gorm:"primary_key;auto_increment"`
	StreamId           uint64
	Stream             *Stream `gorm:"references:StreamId"`
	UpdateFrequency    uint64
	// TODO revisar esta relacion
	RedundantStreams       []*ConfiguredStream `gorm:"references:ConfiguredStreamId"`
	CheckErrors            bool
	NormalUpperThreshold   uint64
	NormalLowerThreshold   uint64
	UnusualNormalThreshold uint64
	UnusualLowerThreshold  uint64
	CalibrationId          uint64
}

func NewConfiguredStream(
	stream *Stream, updateFrequency uint64,
	redundantStreams []*ConfiguredStream, checkErrors bool,
	normalUpperThreshold uint64, normalLowerThreshold uint64,
	unusualNormalThreshold uint64, unusualLowerThreshold uint64,
	calibrationId uint64,
) *ConfiguredStream {
	return &ConfiguredStream{
		Stream:                 stream,
		UpdateFrequency:        updateFrequency,
		RedundantStreams:       redundantStreams,
		CheckErrors:            checkErrors,
		NormalUpperThreshold:   normalUpperThreshold,
		NormalLowerThreshold:   normalLowerThreshold,
		UnusualNormalThreshold: unusualNormalThreshold,
		UnusualLowerThreshold:  unusualLowerThreshold,
		CalibrationId:          calibrationId,
	}
}
