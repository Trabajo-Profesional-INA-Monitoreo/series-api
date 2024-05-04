package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/converters"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	log "github.com/sirupsen/logrus"
)

type ConfigurationService interface {
	GetAllConfigurations() []dtos.AllConfigurations
	GetConfigurationById(id uint64) *dtos.Configuration
	CreateConfiguration(configuration dtos.CreateConfiguration) error
	DeleteConfiguration(id uint64)
	ModifyConfiguration(configuration dtos.Configuration) error
}

type configurationService struct {
	configurationRepository      repositories.ConfigurationRepository
	nodeRepository               repositories.NodeRepository
	configuratedStreamRepository repositories.ManagerConfiguredStreamsRepository
	streamService                StreamCreationService
	metricsRepository            repositories.MetricsRepository
	redundancyRepository         repositories.RedundancyRepository
}

func (c configurationService) deleteFromDbConfiguration(configId uint64) {

	c.configuratedStreamRepository.DeleteByConfig(configId)
	c.nodeRepository.DeleteByConfig(configId)
	c.configurationRepository.DeleteById(configId)

}

func (c configurationService) ModifyConfiguration(configuration dtos.Configuration) error {
	log.Debugf("Updating configuration %v", configuration.Id)
	newConfiguration := converters.ConvertDtoToConfiguration(configuration)
	err := c.configurationRepository.Update(newConfiguration)
	if err != nil {
		return err
	}

	log.Debugf("Updating nodes for configuration %v", configuration.Id)
	var updatedNodesIds []uint64
	for _, node := range configuration.Nodes {
		if node.Id != 0 {
			updatedNodesIds = append(updatedNodesIds, node.Id)
		}
	}
	c.nodeRepository.MarkAsDeletedOldNodes(configuration.Id, updatedNodesIds)

	newNodes := converters.ConvertDtoToNodeModify(*newConfiguration, configuration.Nodes)
	for i, newNode := range newNodes {
		if newNode.NodeId == 0 {
			log.Debugf("Added a new node in configuration %v", configuration.Id)
			nodeId, err := c.nodeRepository.Create(&newNode)
			if err != nil {
				return err
			}
			newNode.NodeId = nodeId
			newNodes[i] = newNode
			continue
		}
		err := c.nodeRepository.Update(&newNode)
		if err != nil {
			return err
		}
	}
	nodes := configuration.Nodes
	var updatedConfigStreamIds []uint64
	for _, node := range nodes {
		for _, configStream := range *node.ConfiguredStreams {
			if configStream.ConfiguredStreamId != 0 {
				updatedConfigStreamIds = append(updatedConfigStreamIds, configStream.ConfiguredStreamId)
			}
		}
	}
	c.configuratedStreamRepository.MarkAsDeletedOldConfiguredStreams(configuration.Id, updatedConfigStreamIds)

	for i, node := range nodes {
		log.Debugf("Updating configured streams for node %v", node.Id)

		for _, configuratedStream := range *node.ConfiguredStreams {

			newConfiguratedStreams := converters.ConvertDtoToConfiguratedStreamModify(*node, configuratedStream, *newConfiguration)

			err := c.streamService.CreateStream(newConfiguratedStreams.StreamId, configuratedStream.StreamType)
			if err != nil {
				return err
			}

			if node.Id == 0 {
				// The stream is in a new node
				newConfiguratedStreams.NodeId = newNodes[i].NodeId
			}

			if configuratedStream.ConfiguredStreamId == 0 {
				log.Debugf("Added a new configured stream in node %v", newConfiguratedStreams.NodeId)
				configuredStreamId, err := c.configuratedStreamRepository.Create(&newConfiguratedStreams)
				if err != nil {
					return err
				}
				configuratedStream.ConfiguredStreamId = configuredStreamId
			} else {
				err = c.configuratedStreamRepository.Update(&newConfiguratedStreams)
				if err != nil {
					return err
				}
			}
			var newMetricsIds []entities.Metric
			if configuratedStream.Metrics != nil {
				for _, metric := range *configuratedStream.Metrics {
					newMetricsIds = append(newMetricsIds, metric)
				}
			}

			c.metricsRepository.DeleteMetricsNotIncludedInNewConfig(configuratedStream.ConfiguredStreamId, newMetricsIds)

			if configuratedStream.Metrics != nil {
				savedMetrics := c.metricsRepository.GetByConfiguredStreamId(configuratedStream.ConfiguredStreamId)
				for _, metric := range *configuratedStream.Metrics {
					// We perform this check to prevent duplicate records on the DB
					found := false
					for i := 0; savedMetrics != nil && i < len(*savedMetrics) && !found; i++ {
						found = (*savedMetrics)[i] == metric
					}
					if found {
						continue
					}
					err = c.metricsRepository.Create(entities.ConfiguredMetric{MetricId: metric, ConfiguredStreamId: configuratedStream.ConfiguredStreamId})
					if err != nil {
						return err
					}
				}
			}

			var newRedundancyIds []uint64
			if configuratedStream.RedundanciesIds != nil {

				for _, id := range *configuratedStream.RedundanciesIds {
					newRedundancyIds = append(newRedundancyIds, id)
				}
			}

			c.redundancyRepository.DeleteRedundanciesNotIncludedInNewConfig(configuratedStream.ConfiguredStreamId, newRedundancyIds)
			if configuratedStream.RedundanciesIds != nil {
				savedRedundancies := c.redundancyRepository.GetByConfiguredStreamId(configuratedStream.ConfiguredStreamId)
				for _, redundancyId := range *configuratedStream.RedundanciesIds {
					// We perform this check to prevent duplicate records on the DB
					found := false
					for i := 0; savedRedundancies != nil && i < len(*savedRedundancies) && !found; i++ {
						found = (*savedRedundancies)[i] == redundancyId
					}
					if found {
						continue
					}
					err = c.redundancyRepository.Create(entities.Redundancy{RedundancyId: redundancyId, ConfiguredStreamId: configuratedStream.ConfiguredStreamId})
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return err

}

func (c configurationService) DeleteConfiguration(id uint64) {
	c.configurationRepository.MarkAsDeleted(id)
	c.nodeRepository.MarkAsDeletedOldNodes(id, nil)
	c.configuratedStreamRepository.MarkAsDeletedOldConfiguredStreams(id, nil)
	go c.deleteFromDbConfiguration(id)
}

func (c configurationService) CreateConfiguration(configuration dtos.CreateConfiguration) error {
	log.Debugf("Creating configuration %v", configuration.Name)
	newConfiguration := converters.ConvertDtoToCreateConfiguration(configuration)
	err := c.configurationRepository.Create(newConfiguration)
	if err != nil {
		return err
	}
	newNodes := converters.ConvertDtoToNode(*newConfiguration, configuration.Nodes)

	log.Debugf("Creating nodes for configuration %v", configuration.Name)
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
		log.Debugf("Creating configured streams for node %v", node.Name)
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

func (c configurationService) GetConfigurationById(id uint64) *dtos.Configuration {
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
		if node.ConfiguredStreams == nil {
			log.Warnf("Saved node without configured streams, skipping...")
			continue
		}
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
		streamService:                NewStreamService(repositories, client),
		metricsRepository:            repositories.MetricsRepository,
		redundancyRepository:         repositories.RedundancyRepository,
	}
}
