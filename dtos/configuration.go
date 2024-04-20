package dtos

type Configuration struct {
	Name              string `binding:"required" ,gorm:"column:name"`
	Id                uint64 `gorm:"column:configuration_id"`
	SendNotifications bool
	Nodes             []*Node `gorm:"-"`
}

type AllConfigurations struct {
	Name string `gorm:"column:name"`
	Id   uint64 `gorm:"column:configuration_id"`
}

type CreateConfiguration struct {
	Name              string `binding:"required" ,gorm:"column:name"`
	SendNotifications bool
	Nodes             *[]CreateNode
}
