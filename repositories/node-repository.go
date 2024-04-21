package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type NodeRepository interface {
	Create(node *entities.Node) (uint64, error)
	Update(node *entities.Node) error
	GetNodesById(id uint64) []*dtos.Node
	GetStreamsPerNodeById(formatUint string) []*dtos.StreamsPerNode
	MarkAsDeletedOldNodes(id uint64, ids []uint64)
}

type nodeRepository struct {
	connection *gorm.DB
}

func (n nodeRepository) MarkAsDeletedOldNodes(configId uint64, newNodeIds []uint64) {
	n.connection.Model(
		&entities.Node{},
	).Where(
		"nodes.configuration_id = ?", configId,
	).Where(
		"nodes.node_id NOT IN ?", newNodeIds,
	).Update("deleted", true)
}

func (n nodeRepository) GetStreamsPerNodeById(configId string) []*dtos.StreamsPerNode {
	var nodes []*dtos.StreamsPerNode

	tx := n.connection.Model(
		&entities.Node{},
	).Select(
		"nodes.name as name",
		"nodes.node_id as node_id",
		"count(distinct(configured_streams.stream_id)) as streams_count",
	).Joins(
		"JOIN configured_streams ON configured_streams.node_id = nodes.node_id",
	).Where(
		"nodes.configuration_id = ?", configId,
	).Group(
		"nodes.name, nodes.node_id",
	).Scan(&nodes)

	if tx.Error != nil {
		log.Errorf("Error executing GetNodes query: %v", tx.Error)
	}

	log.Debugf("Get nodes query result: %v", nodes)
	return nodes
}

func (n nodeRepository) Update(node *entities.Node) error {
	result := n.connection.Save(&node)
	return result.Error
}

func (n nodeRepository) Create(node *entities.Node) (uint64, error) {
	result := n.connection.Create(&node)
	return node.NodeId, result.Error
}

func (n nodeRepository) GetNodesById(id uint64) []*dtos.Node {
	var nodes []*dtos.Node

	result := n.connection.Model(
		&entities.Node{},
	).Select(
		"nodes.name as name, nodes.node_id as node_id",
	).Where("configuration_id = ?", id).Where("deleted = false").Scan(&nodes)

	if result.RowsAffected == 0 {
		return nil
	}

	log.Debugf("Get node query result: %v", nodes)
	return nodes
}

func NewNodeRepository(connection *gorm.DB) NodeRepository {
	return &nodeRepository{connection}
}
