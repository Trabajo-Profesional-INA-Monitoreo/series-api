package dtos

type Configuration struct {
	Name              string  `binding:"required,min=1" ,gorm:"column:name" json:"Name"`
	Id                uint64  `gorm:"column:configuration_id" json:"Id"`
	SendNotifications bool    `json:"SendNotifications"`
	Nodes             []*Node `binding:"required,min=1,dive" gorm:"-" json:"Nodes"`
}

type AllConfigurations struct {
	Name string `gorm:"column:name" json:"Name"`
	Id   uint64 `gorm:"column:configuration_id" json:"Id"`
}

type CreateConfiguration struct {
	Name              string        `binding:"required,min=1" ,gorm:"column:name" json:"Name"`
	SendNotifications bool          `json:"SendNotifications"`
	Nodes             *[]CreateNode `binding:"required,min=1,dive" json:"Nodes"`
}
