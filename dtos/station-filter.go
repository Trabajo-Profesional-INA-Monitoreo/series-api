package dtos

type StationFilter struct {
	Name string `gorm:"column:name" json:"Name"`
	Id   uint64 `gorm:"column:id" json:"Id"`
}
