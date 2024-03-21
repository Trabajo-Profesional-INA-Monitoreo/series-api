package entities

type Variable struct {
	VariableId uint64 `gorm:"primary_key"`
	Name       string `gorm:"type:varchar(100)"`
}
