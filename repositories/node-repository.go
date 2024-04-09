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
	GetNodesById(id string) []*dtos.Node
}

type nodeRepository struct {
	connection *gorm.DB
}

func (n nodeRepository) Update(node *entities.Node) error {
	result := n.connection.Save(&node)
	return result.Error
}

func (n nodeRepository) Create(node *entities.Node) (uint64, error) {
	result := n.connection.Create(&node)
	return node.NodeId, result.Error
}

func (n nodeRepository) GetNodesById(id string) []*dtos.Node {
	var nodes []*dtos.Node

	result := n.connection.Model(
		&entities.Node{},
	).Select(
		"nodes.name as name, nodes.node_id as node_id",
	).Where("configuration_id = ?", id).Scan(&nodes)

	if result.RowsAffected == 0 {
		return nil
	}

	log.Debugf("Get node query result: %v", nodes)
	return nodes
}

func NewNodeRepository(connection *gorm.DB) NodeRepository {
	return &nodeRepository{connection}
}
