package entities

type Stream struct {
	StreamId   uint64 `gorm:"primary_key"`
	StationId  uint64
	Station    *Station `gorm:"references:StationId"`
	NetworkId  uint64
	Network    *Network `gorm:"references:NetworkId"`
	Variable   string   `gorm:"type:varchar(100)"`
	VariableId uint64
	ProcId     uint64
	StreamType StreamType
}

func NewStream(streamId uint64, station *Station, network *Network, variable string, variableId uint64, procId uint64, streamType StreamType) *Stream {
	return &Stream{
		StreamId:   streamId,
		Station:    station,
		Network:    network,
		Variable:   variable,
		VariableId: variableId,
		ProcId:     procId,
		StreamType: streamType,
	}
}
