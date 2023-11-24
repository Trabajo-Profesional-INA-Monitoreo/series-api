package entities

type Station struct {
	StationId uint64 `gorm:"primary_key"`
	Name      string `gorm:"type:varchar(100)"`
}

func NewStation(stationId uint64, name string) *Station {
	return &Station{StationId: stationId, Name: name}
}
