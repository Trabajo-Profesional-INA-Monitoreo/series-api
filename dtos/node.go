package dtos

type Node struct {
	Name              string               `binding:"required,min=1" gorm:"column:name" json:"Name"`
	Id                uint64               `gorm:"column:node_id" json:"Id"`
	ConfiguredStreams *[]*ConfiguredStream `binding:"required,min=1,dive" gorm:"-" json:"ConfiguredStreams"`
	MainStreamId      *uint64              `binding:"omitempty,min=1" json:"MainStreamId"`
}

type CreateNode struct {
	Name              string                    `binding:"required,min=1" gorm:"column:name" json:"Name"`
	ConfiguredStreams *[]ConfiguredStreamCreate `binding:"required,min=1,dive" gorm:"-" json:"ConfiguredStreams"`
	MainStreamId      *uint64                   `binding:"omitempty,min=1" json:"MainStreamId"`
}
