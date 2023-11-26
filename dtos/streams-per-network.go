package dtos

type StreamsPerNetworkResponse struct {
	Networks []StreamsPerNetwork
}

type StreamsPerNetwork struct {
	NetworkName  string `gorm:"column:networkname"`
	NetworkId    string `gorm:"column:networkid"`
	StreamsCount int    `gorm:"column:streamscount"`
}
