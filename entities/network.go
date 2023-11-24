package entities

type Network struct {
	NetworkId uint64 `gorm:"primary_key"`
	Name      string `gorm:"type:varchar(100)"`
}

func NewNetwork(networkId uint64, name string) *Network {
	return &Network{NetworkId: networkId, Name: name}
}
