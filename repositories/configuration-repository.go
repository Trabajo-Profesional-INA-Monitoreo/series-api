package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ConfigurationRepository interface {
	Save(configuration *entities.Configuration) error
	GetAllConfigurations() []dtos.AllConfigurations
	GetConfigurationById(id string) *dtos.Configuration
	Delete(id string)
}

type configurationRepository struct {
	connection *gorm.DB
}

func (c configurationRepository) Delete(id string) {
	c.connection.Where("name = ?", id).Delete(&dtos.Configuration{})
}

func (c configurationRepository) GetConfigurationById(id string) *dtos.Configuration {
	var configuration dtos.Configuration

	result := c.connection.Model(
		&entities.Configuration{},
	).Select(
		"configurations.name as name",
	).Where("name = ?", id).Scan(&configuration)

	if result.RowsAffected == 0 {
		return nil
	}

	log.Debugf("Get configurations query result: %v", configuration)
	return &configuration
}

func (c configurationRepository) GetAllConfigurations() []dtos.AllConfigurations {
	var configurations []dtos.AllConfigurations

	c.connection.Model(
		&entities.Configuration{},
	).Select(
		"configurations.name as name",
	).Scan(&configurations)

	log.Debugf("Get configurations query result: %v", configurations)
	return configurations
}

func (c configurationRepository) Save(configuration *entities.Configuration) error {
	result := c.connection.Create(&configuration)
	return result.Error
}

func NewConfigurationRepository(connection *gorm.DB) ConfigurationRepository {
	return &configurationRepository{connection}
}