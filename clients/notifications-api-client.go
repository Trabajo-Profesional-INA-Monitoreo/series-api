package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	exceptions "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type NotificationsAPiClient interface {
	SendNotification(message string) error
}

type notificationsAPiClient struct {
	baseUrl string
}

func (n notificationsAPiClient) SendNotification(message string) error {
	reqUrl := fmt.Sprintf("%v/notificacion", n.baseUrl)
	log.Debugf("Send notification: %v", reqUrl)

	values := map[string]string{"Message": message}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Errorf("Error making json: %v", err)
		return err
	}

	res, err := http.Post(reqUrl, "application/json", bytes.NewBuffer(json_data))

	if err != nil {
		log.Errorf("Error making request: %v", err)
		return err
	}

	defer closeReaderAndPrintError(res.Body)

	if res.StatusCode != 201 {
		err := fmt.Errorf("notifications api response error: got %v", res.StatusCode)
		log.Errorf(err.Error())
		return errors.Join(exceptions.MapCodeToError(res.StatusCode), err)
	}

	return nil
}

func NewNotificationsAPiClientImpl(apiConfig *config.ServiceConfigurationData) NotificationsAPiClient {
	return &notificationsAPiClient{baseUrl: apiConfig.NotificationApiBaseUrl}
}
