package dtos

type ProcedureFilter struct {
	Name string `gorm:"column:name"`
	Id   uint64 `gorm:"column:id"`
}
