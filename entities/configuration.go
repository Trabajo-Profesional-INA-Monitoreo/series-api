package entities

import "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"

type Configuration struct {
	Name string `gorm:"primary_key;type:varchar(100)"`
}

func NewConfiguration(configuration dtos.Configuration) *Configuration {
	return &Configuration{Name: configuration.Name}
}
