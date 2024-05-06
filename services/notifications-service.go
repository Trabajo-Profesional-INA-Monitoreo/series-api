package services

import "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"

type NotificationsService interface {
	SendDailyNotification()
}

type notificationsService struct {
	client clients.NotificationsAPiClient
}

func (n notificationsService) SendDailyNotification() {
	n.client.SendNotification("Increiblemente, esto viene de un job. No te la veias venir")
}

func NewNotificationsService(client clients.NotificationsAPiClient) NotificationsService {
	return &notificationsService{client: client}
}
