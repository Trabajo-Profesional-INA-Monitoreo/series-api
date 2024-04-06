package dtos

type VariableFilter struct {
	Name string `gorm:"column:name"`
	Id   uint64 `gorm:"column:id"`
}
