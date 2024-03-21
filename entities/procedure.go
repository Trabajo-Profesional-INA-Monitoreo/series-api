package entities

type Procedure struct {
	ProcedureId uint64 `gorm:"primary_key"`
	Name        string `gorm:"type:varchar(100)"`
}

func NewProcedure(procId uint64, name string) *Procedure {
	return &Procedure{ProcedureId: procId, Name: name}
}
