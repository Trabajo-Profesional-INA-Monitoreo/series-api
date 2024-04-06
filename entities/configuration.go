package entities

type Configuration struct {
	Name            string `gorm:"type:varchar(100)"`
	ConfigurationId uint64 `gorm:"primary_key;auto_increment"`
}
