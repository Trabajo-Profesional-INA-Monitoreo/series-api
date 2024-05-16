package entities

type Unit struct {
	UnitId uint64 `gorm:"primary_key"`
	Name   string `gorm:"type:varchar(100)"`
}
