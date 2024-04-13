package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/converters"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
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
	metricsRepository            repositories.MetricsRepository
	redundancyRepository         repositories.RedundancyRepository
}

func (c configurationService) ModifyConfiguration(configuration dtos.Configuration) error {

	newConfiguration := converters.ConvertDtoToConfiguration(configuration)
	err := c.configurationRepository.Update(newConfiguration)
	if err != nil {
		return err
	}
	newNodes := converters.ConvertDtoToNodeModify(*newConfiguration, configuration.Nodes)
	for _, newNode := range newNodes {
		err := c.nodeRepository.Update(&newNode)
		if err != nil {
			return err
		}
	}

	nodes := configuration.Nodes
	for _, node := range nodes {

		for _, configuratedStream := range *node.ConfiguredStreams {

			newConfiguratedStreams := converters.ConvertDtoToConfiguratedStreamModify(*node, configuratedStream, *newConfiguration)

			err := c.streamService.CreateStream(newConfiguratedStreams.StreamId, configuratedStream.StreamType)
			if err != nil {
				return err
			}

			err = c.configuratedStreamRepository.Update(&newConfiguratedStreams)
			if err != nil {
				return err
			}
		}
	}
	return err

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

	var nodeIds []uint64

	for _, newNode := range newNodes {
		nodeId, err := c.nodeRepository.Create(&newNode)
		nodeIds = append(nodeIds, nodeId)
		if err != nil {
			return err
		}
	}

	nodes := *configuration.Nodes
	for index, node := range nodes {

		for _, configuratedStream := range *node.ConfiguredStreams {

			newConfiguratedStreams := converters.ConvertDtoToConfiguratedStream(nodeIds[index], &configuratedStream, *newConfiguration)

			err := c.streamService.CreateStream(newConfiguratedStreams.StreamId, configuratedStream.StreamType)
			if err != nil {
				return err
			}

			var configuredStreamId uint64
			configuredStreamId, err = c.configuratedStreamRepository.Create(&newConfiguratedStreams)
			if err != nil {
				return err
			}

			for _, metric := range configuratedStream.Metrics {
				err = c.metricsRepository.Create(entities.ConfiguredMetric{MetricId: metric, ConfiguredStreamId: configuredStreamId})
				if err != nil {
					return err
				}
			}

			for _, redundancyId := range configuratedStream.RedundanciesIds {
				err = c.redundancyRepository.Create(entities.Redundancy{RedundancyId: redundancyId, ConfiguredStreamId: configuredStreamId})
				if err != nil {
					return err
				}
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
	nodes = c.nodeRepository.GetNodesById(id)
	if nodes == nil {
		return nil
	}
	configuration.Nodes = nodes

	for _, node := range nodes {
		node.ConfiguredStreams = c.configuratedStreamRepository.FindConfiguredStreamsByNodeId(node.Id, id)
		for _, configuredStream := range *node.ConfiguredStreams {
			configuredStream.Metrics = c.metricsRepository.GetByConfiguredStreamId(configuredStream.ConfiguredStreamId)
			configuredStream.RedundanciesIds = c.redundancyRepository.GetByConfiguredStreamId(configuredStream.ConfiguredStreamId)
		}
	}

	return configuration
}

func NewConfigurationService(repositories *config.Repositories, client clients.InaAPiClient) ConfigurationService {
	return &configurationService{configurationRepository: repositories.ConfigurationRepository,
		nodeRepository:               repositories.NodeRepository,
		configuratedStreamRepository: repositories.ConfiguredStreamRepository,
		streamService:                NewStreamService(repositories.StreamsRepository, client, repositories.ConfiguredStreamRepository),
		metricsRepository:            repositories.MetricsRepository,
		redundancyRepository:         repositories.RedundancyRepository,
	}
}
