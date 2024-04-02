package dtos

type StationFilter struct {
	Name string `gorm:"column:name"`
	Id   uint64 `gorm:"column:id"`
}
