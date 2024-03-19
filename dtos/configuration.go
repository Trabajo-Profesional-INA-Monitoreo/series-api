package dtos

type Configuration struct {
	Name string `binding:"required" ,gorm:"column:name"`
	Id   uint64 `gorm:"column:id"`
}

type AllConfigurations struct {
	Name string `gorm:"column:name"`
	Id   uint64 `gorm:"column:id"`
}
