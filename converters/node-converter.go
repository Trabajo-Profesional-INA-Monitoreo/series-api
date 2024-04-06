package converters

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
)

func ConvertDtoToNode(configuration entities.Configuration, nodes *[]dtos.Node) []entities.Node {
	var newNodes []entities.Node

	for _, node := range *nodes {
		newNodes = append(newNodes, entities.Node{ConfigurationId: configuration.ConfigurationId, Name: node.Name, NodeId: node.Id})
	}

	return newNodes
}

func ConvertDtoToNodeModify(configuration entities.Configuration, nodes []*dtos.Node) []entities.Node {
	var newNodes []entities.Node

	for _, node := range nodes {
		newNodes = append(newNodes, entities.Node{ConfigurationId: configuration.ConfigurationId, Name: node.Name, NodeId: node.Id})
	}

	return newNodes
}
