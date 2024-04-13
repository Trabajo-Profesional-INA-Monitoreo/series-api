package entities

type ConfiguredStream struct {
	ConfiguredStreamId     uint64 `gorm:"primary_key;auto_increment"`
	StreamId               uint64
	Stream                 *Stream `gorm:"references:StreamId"`
	UpdateFrequency        float64
	RedundantStreams       []Redundancy
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
