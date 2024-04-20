package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ConfigurationRepository interface {
	Create(configuration *entities.Configuration) error
	GetAllConfigurations() []dtos.AllConfigurations
	GetConfigurationById(id string) *dtos.Configuration
	Delete(id string)
	Update(configuration *entities.Configuration) error
}

type configurationRepository struct {
	connection *gorm.DB
}

func (c configurationRepository) Update(configuration *entities.Configuration) error {
	result := c.connection.Save(&configuration)
	return result.Error
}

func (c configurationRepository) Delete(id string) {
	c.connection.Model(&entities.Configuration{}).Where("configuration_id = ?", id).Update("deleted", true)
}

func (c configurationRepository) GetConfigurationById(id string) *dtos.Configuration {
	var configuration dtos.Configuration

	result := c.connection.Model(
		&entities.Configuration{},
	).Select(
		"configurations.name as name",
		"configurations.configuration_id as configuration_id",
		"configurations.send_notifications as send_notifications",
	).Where("configuration_id = ? AND deleted = false", id).Scan(&configuration)

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
		"configurations.name as name, configurations.configuration_id as configuration_id",
	).Where("deleted = false").Scan(&configurations)

	log.Debugf("Get configurations query result: %v", configurations)
	return configurations
}

func (c configurationRepository) Create(configuration *entities.Configuration) error {
	result := c.connection.Create(&configuration)
	return result.Error
}

func NewConfigurationRepository(connection *gorm.DB) ConfigurationRepository {
	return &configurationRepository{connection}
}
