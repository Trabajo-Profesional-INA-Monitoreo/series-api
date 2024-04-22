package dtos

type Configuration struct {
	Name              string `binding:"required,min=1" ,gorm:"column:name"`
	Id                uint64 `gorm:"column:configuration_id"`
	SendNotifications bool
	Nodes             []*Node `binding:"required,min=1,dive" gorm:"-"`
}

type AllConfigurations struct {
	Name string `gorm:"column:name"`
	Id   uint64 `gorm:"column:configuration_id"`
}

type CreateConfiguration struct {
	Name              string `binding:"required,min=1" ,gorm:"column:name"`
	SendNotifications bool
	Nodes             *[]CreateNode `binding:"required,min=1,dive"`
}
