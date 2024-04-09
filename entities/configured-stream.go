package entities

type ConfiguredStream struct {
	ConfiguredStreamId     uint64 `gorm:"primary_key;auto_increment"`
	StreamId               uint64
	Stream                 *Stream `gorm:"references:StreamId"`
	UpdateFrequency        float64
	RedundantStreams       []ConfiguredStream `gorm:"many2many:redundancies;"`
	CheckErrors            bool
	NormalUpperThreshold   float64
	NormalLowerThreshold   float64
	UnusualNormalThreshold float64
	UnusualLowerThreshold  float64
	CalibrationId          uint64
	Metrics                []ConfiguredMetric
	NodeId                 uint64
	Node                   *Node `gorm:"references:NodeId"`
	ConfigurationId        uint64
	Configuration          *Configuration `gorm:"references:ConfigurationId"`
}

func NewConfiguredStream(
	stream *Stream, updateFrequency float64,
	redundantStreams []ConfiguredStream, checkErrors bool,
	normalUpperThreshold float64, normalLowerThreshold float64,
	calibrationId uint64, Metrics []ConfiguredMetric,
) *ConfiguredStream {
	return &ConfiguredStream{
		Stream:               stream,
		UpdateFrequency:      updateFrequency,
		RedundantStreams:     redundantStreams,
		CheckErrors:          checkErrors,
		NormalUpperThreshold: normalUpperThreshold,
		NormalLowerThreshold: normalLowerThreshold,
		CalibrationId:        calibrationId,
		Metrics:              Metrics,
	}
}
