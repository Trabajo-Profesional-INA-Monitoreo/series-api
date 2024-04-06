package dtos

type Node struct {
	Name              string              `gorm:"column:name"`
	Id                uint64              `gorm:"column:node_id"`
	ConfiguredStreams *[]ConfiguredStream `gorm:"-"`
}
