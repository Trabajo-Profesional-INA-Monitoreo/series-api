package entities

type Node struct {
	Name            string `gorm:"type:varchar(100)"`
	NodeId          uint64 `gorm:"primary_key;auto_increment"`
	ConfigurationId uint64
	Configuration   *Configuration `gorm:"references:ConfigurationId"`
	Deleted         bool
	MainStreamId    *uint64
}
