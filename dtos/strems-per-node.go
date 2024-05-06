package dtos

type StreamsPerNodeResponse struct {
	Nodes []*StreamsPerNode `json:"Nodes"`
}

type StreamsPerNode struct {
	NodeName     string `gorm:"column:name" json:"NodeName"`
	NodeId       string `gorm:"column:node_id" json:"NodeId"`
	StreamsCount int    `gorm:"column:streams_count" json:"StreamsCount"`
	ErrorCount   int    `json:"ErrorCount"`
}

type ErrorsOfNodes struct {
	NodeId     string `gorm:"column:node_id" json:"NodeId"`
	ErrorCount int    `gorm:"column:error_count" json:"ErrorCount"`
}
