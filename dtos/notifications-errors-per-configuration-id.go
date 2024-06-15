package dtos

type NotificationsErrorsCountPerConfigurationId struct {
	ConfigurationId string `gorm:"column:configuration_id"`
	Name            string `gorm:"column:name"`
	Total           int    `gorm:"column:total"`
}
