package dtos

type Node struct {
	Name              string               `binding:"required,min=1" gorm:"column:name"`
	Id                uint64               `gorm:"column:node_id"`
	ConfiguredStreams *[]*ConfiguredStream `binding:"required,min=1,dive" gorm:"-"`
}

type CreateNode struct {
	Name              string                    `binding:"required,min=1" gorm:"column:name"`
	ConfiguredStreams *[]ConfiguredStreamCreate `binding:"required,min=1,dive" gorm:"-"`
}
