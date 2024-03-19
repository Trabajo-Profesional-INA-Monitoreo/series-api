package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/converters"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
)

type (
	ConfigurationService interface {
		GetAllConfigurations() []dtos.AllConfigurations
		GetConfigurationById(id string) *dtos.Configuration
		CreateConfiguration(configuration dtos.Configuration) error
		DeleteConfiguration(id string)
		ModifyConfiguration(configuration dtos.Configuration) error
	}
)

type configurationService struct {
	repository repositories.ConfigurationRepository
}

func (c configurationService) ModifyConfiguration(configuration dtos.Configuration) error {
	newConfiguration := converters.ConvertDtoToConfiguration(configuration)
	return c.repository.Update(newConfiguration)
}

func (c configurationService) DeleteConfiguration(id string) {
	c.repository.Delete(id)
}

func (c configurationService) CreateConfiguration(configuration dtos.Configuration) error {
	newConfiguration := converters.ConvertDtoToConfiguration(configuration)
	return c.repository.Create(newConfiguration)
}

func (c configurationService) GetAllConfigurations() []dtos.AllConfigurations {
	return c.repository.GetAllConfigurations()
}

func (c configurationService) GetConfigurationById(id string) *dtos.Configuration {
	return c.repository.GetConfigurationById(id)
}

func NewConfigurationService(repository repositories.ConfigurationRepository) ConfigurationService {
	return &configurationService{repository}
}
