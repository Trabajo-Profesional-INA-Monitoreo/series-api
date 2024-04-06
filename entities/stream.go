package entities

type Stream struct {
	StreamId  uint64 `gorm:"primary_key"`
	StationId uint64
	Station   *Station `gorm:"references:StationId"`
	//NetworkId   uint64
	//Network     *Network `gorm:"references:NetworkId"`
	VariableId  uint64
	Variable    *Variable `gorm:"references:VariableId"`
	ProcedureId uint64
	Procedure   *Procedure `gorm:"references:ProcedureId"`
	UnitId      uint64
	Unit        *Unit `gorm:"references:UnitId"`
	StreamType  StreamType
}

func NewStream(streamId uint64, station *Station, network *Network, variable *Variable, procedure *Procedure, unit *Unit, streamType StreamType) *Stream {
	return &Stream{
		StreamId: streamId,
		Station:  station,
		//Network:    network,
		Variable:   variable,
		Procedure:  procedure,
		Unit:       unit,
		StreamType: streamType,
	}
}

func (s Stream) IsForecasted() bool {
	return s.StreamType == Forecasted
}

func (s Stream) IsObserved() bool {
	return s.StreamType == Observed
}
