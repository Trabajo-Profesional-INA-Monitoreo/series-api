package entities

type Configuration struct {
	Name string `gorm:"type:varchar(100)"`
	Id   uint64 `gorm:"primary_key"`
}

//func NewConfiguration(configuration dtos.Configuration) *Configuration {
//	return &Configuration{Name: configuration.Name, Id: configuration.Id}
//}
