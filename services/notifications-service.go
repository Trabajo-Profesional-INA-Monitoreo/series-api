package services

import (
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"time"
)

type NotificationsService interface {
	SendDailyNotification()
}

type notificationsService struct {
	client                  clients.NotificationsAPiClient
	notificationsRepository repositories.NotificationsRepository
}

func (n notificationsService) SendDailyNotification() {
	yesterday := time.Now().Add(-24 * time.Hour)
	totalErrors := n.notificationsRepository.GetTotalErrors(yesterday)
	errorsByConfigurationId := n.notificationsRepository.GetErrorsByConfigurationId(yesterday)
	errorsByConfigurationIdAndErrorType := n.notificationsRepository.GetErrorsByConfigurationIdAndErrorType(yesterday)
	message := n.createMessage(yesterday, totalErrors, errorsByConfigurationId, errorsByConfigurationIdAndErrorType)
	n.client.SendNotification(message)
}

func (n notificationsService) createMessage(date time.Time, totalErrors int64, errorsByConfigurationId []*dtos.NotificationsErrorsCountPerConfigurationId, errorsByConfigurationIdAndErrorType []*dtos.NotificationsErrorsCountPerType) string {
	mensaje := fmt.Sprintf("Resumen del día %v/%v/%v \n\n", date.Day(), int(date.Month()), date.Year())

	if totalErrors == 0 {
		mensaje += fmt.Sprintf("Se detectaron %v errores totales.\n", totalErrors)
		return mensaje
	}
	mensaje += fmt.Sprintf("Se detectaron %v errores totales:\n\n", totalErrors)

	for _, configuration := range errorsByConfigurationId {
		mensaje += fmt.Sprintf("* En la configuración con id %v \"%v\", se detectaron %v errores.\n\n", configuration.ConfigurationId, configuration.Name, configuration.Total)
		mensaje += fmt.Sprintf("De los cuales:\n")
		for _, errorType := range errorsByConfigurationIdAndErrorType {
			if configuration.ConfigurationId == errorType.ConfigurationId {
				mensaje += fmt.Sprintf("- %v son de tipo %v.\n", errorType.Total, errorType.ErrorTypeId.Translate())
			}
		}
		mensaje += fmt.Sprintf("\n")
	}

	return mensaje
}

func NewNotificationsService(client clients.NotificationsAPiClient, repositories *config.Repositories) NotificationsService {
	return &notificationsService{client: client, notificationsRepository: repositories.NotificationsRepository}
}
