package entities

type ConfiguredStream struct {
	ConfiguredStreamId     uint64 `gorm:"primary_key;auto_increment"`
	ConfigurationId        uint64
	Configuration          *Configuration `gorm:"references:Id"`
	StreamId               uint64
	Stream                 *Stream `gorm:"references:StreamId"`
	UpdateFrequency        float64
	RedundantStreams       []ConfiguredStream `gorm:"many2many:redundancies;"`
	CheckErrors            bool
	NormalUpperThreshold   uint64
	NormalLowerThreshold   uint64
	UnusualNormalThreshold uint64
	UnusualLowerThreshold  uint64
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
	normalUpperThreshold uint64, normalLowerThreshold uint64,
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
