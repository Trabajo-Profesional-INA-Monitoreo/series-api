package entities

type Stream struct {
	StreamId    uint64 `gorm:"primary_key"`
	StationId   uint64
	Station     *Station `gorm:"references:StationId"`
	VariableId  uint64
	Variable    *Variable `gorm:"references:VariableId"`
	ProcedureId uint64
	Procedure   *Procedure `gorm:"references:ProcedureId"`
	UnitId      uint64
	Unit        *Unit `gorm:"references:UnitId"`
	StreamType  StreamType
}

func (s Stream) IsForecasted() bool {
	return s.StreamType == Forecasted
}

func (s Stream) IsObserved() bool {
	return s.StreamType == Observed
}
