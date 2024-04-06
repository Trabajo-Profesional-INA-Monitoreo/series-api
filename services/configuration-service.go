package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/converters"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
)

type (
	ConfigurationService interface {
		GetAllConfigurations() []dtos.AllConfigurations
		GetConfigurationById(id string) *dtos.Configuration
		CreateConfiguration(configuration dtos.CreateConfiguration) error
		DeleteConfiguration(id string)
		ModifyConfiguration(configuration dtos.Configuration) error
	}
)

type configurationService struct {
	configurationRepository      repositories.ConfigurationRepository
	nodeRepository               repositories.NodeRepository
	configuratedStreamRepository repositories.ConfiguredStreamsRepository
	streamService                StreamService
}

func (c configurationService) ModifyConfiguration(configuration dtos.Configuration) error {
	newConfiguration := converters.ConvertDtoToConfiguration(configuration)
	return c.configurationRepository.Update(newConfiguration)
}

func (c configurationService) DeleteConfiguration(id string) {
	c.configurationRepository.Delete(id)
}

func (c configurationService) CreateConfiguration(configuration dtos.CreateConfiguration) error {
	newConfiguration := converters.ConvertDtoToCreateConfiguration(configuration)
	err := c.configurationRepository.Create(newConfiguration)
	if err != nil {
		return err
	}
	newNodes := converters.ConvertDtoToNode(*newConfiguration, configuration.Nodes)
	for _, newNode := range newNodes {
		err := c.nodeRepository.Create(&newNode)
		if err != nil {
			return err
		}
	}

	nodes := *configuration.Nodes
	for _, node := range nodes {

		for _, configuratedStream := range *node.ConfiguredStreams {

			newConfiguratedStreams := converters.ConvertDtoToConfiguratedStream(node, &configuratedStream, *newConfiguration)

			err := c.streamService.CreateStream(newConfiguratedStreams.StreamId, configuratedStream.StreamType)
			if err != nil {
				return err
			}

			err = c.configuratedStreamRepository.Create(&newConfiguratedStreams)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (c configurationService) GetAllConfigurations() []dtos.AllConfigurations {
	return c.configurationRepository.GetAllConfigurations()
}

func (c configurationService) GetConfigurationById(id string) *dtos.Configuration {
	var configuration *dtos.Configuration
	var nodes []*dtos.Node

	configuration = c.configurationRepository.GetConfigurationById(id)
	if configuration == nil {
		return nil
	}
	nodes = c.configurationRepository.GetNodesById(id)
	if nodes == nil {
		return nil
	}
	configuration.Nodes = nodes

	for _, node := range nodes {
		node.ConfiguredStreams = c.configuratedStreamRepository.FindConfiguredStreamsByNodeId(node.Id, id)
	}

	return configuration
}

func NewConfigurationService(repositories *config.Repositories, client clients.InaAPiClient) ConfigurationService {
	return &configurationService{configurationRepository: repositories.ConfigurationRepository, nodeRepository: repositories.NodeRepository, configuratedStreamRepository: repositories.ConfiguredStreamRepository, streamService: NewStreamService(repositories.StreamsRepository, client, repositories.ConfiguredStreamRepository)}
}
