package entities

type Procedure struct {
	ProcedureId uint64 `gorm:"primary_key"`
	Name        string `gorm:"type:varchar(100)"`
}
