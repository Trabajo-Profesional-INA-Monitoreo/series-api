package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type NodeRepository interface {
	Create(node *entities.Node) (uint64, error)
	Update(node *entities.Node) error
	GetNodesById(id uint64) []*dtos.Node
	GetStreamsPerNodeById(formatUint uint64) []*dtos.StreamsPerNode
	MarkAsDeletedOldNodes(id uint64, ids []uint64)
	DeleteByConfig(configId uint64)
	GetErrorsOfNodes(configId uint64, timeStart time.Time, timeEnd time.Time) []dtos.ErrorsOfNodes
}

type nodeRepository struct {
	connection *gorm.DB
}

func (n nodeRepository) DeleteByConfig(configId uint64) {
	n.connection.Where("nodes.configuration_id = ?", configId).Where("nodes.deleted = true").Delete(&entities.Node{})
}

func (n nodeRepository) MarkAsDeletedOldNodes(configId uint64, newNodeIds []uint64) {
	tx := n.connection.Model(
		&entities.Node{},
	).Where(
		"nodes.configuration_id = ?", configId,
	)
	if newNodeIds != nil && len(newNodeIds) != 0 {
		tx = tx.Where(
			"nodes.node_id NOT IN ?", newNodeIds,
		)
	}
	tx.Update("deleted", true)
}

func (n nodeRepository) GetStreamsPerNodeById(configId uint64) []*dtos.StreamsPerNode {
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
	).Where(
		"configured_streams.deleted = false",
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

func (n nodeRepository) GetErrorsOfNodes(configId uint64, timeStart time.Time, timeEnd time.Time) []dtos.ErrorsOfNodes {
	var errorsPerNode []dtos.ErrorsOfNodes

	tx := n.connection.Model(
		&entities.ConfiguredStream{},
	).Select(
		"configured_streams.node_id as node_id",
		"count(detected_errors.error_id) as error_count",
	).Joins(
		"JOIN configured_streams_errors ON configured_streams.configured_stream_id = configured_streams_errors.configured_stream_configured_stream_id",
	).Joins(
		"JOIN detected_errors ON detected_errors.error_id = configured_streams_errors.detected_error_error_id ",
	).Where(
		"configured_streams.configuration_id = ?", configId,
	).Where(
		"detected_errors.detected_date >= ? AND detected_errors.detected_date <= ?", timeStart, timeEnd,
	).Where(
		"configured_streams.deleted = false",
	).Group("configured_streams.node_id").Scan(&errorsPerNode)

	if tx.Error != nil {
		log.Errorf("Error executing GetErrorsOfStations query: %v", tx.Error)
	}

	return errorsPerNode
}

func NewNodeRepository(connection *gorm.DB) NodeRepository {
	return &nodeRepository{connection}
}
