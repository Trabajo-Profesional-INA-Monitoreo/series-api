package dtos

type Configuration struct {
	Name string `binding:"required" ,gorm:"column:name"`
}

type AllConfigurations struct {
	Name string `gorm:"column:name"`
}
