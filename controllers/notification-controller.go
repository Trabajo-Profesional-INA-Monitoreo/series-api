package controllers

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NotificationController interface {
	CreateNotification(context *gin.Context)
}

type notificationController struct {
	client clients.NotificationsAPiClient
}

// CreateNotification godoc
//
//			@Summary		Endpoint para crear una notificacion
//			@Tags           Notificacion
//			@Produce		json
//	   		@Param          notification  body  dtos.Notification    true    "Notification"
//			@Success		201
//		    @Failure        400 {object} dtos.ErrorResponse
//			@Router			/notificacion [post]
func (n notificationController) CreateNotification(ctx *gin.Context) {
	var notification dtos.Notification

	if err := ctx.ShouldBindJSON(&notification); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(err))
		return
	}

	err := n.client.SendNotification(notification.Message)

	if err != nil {
		ctx.JSON(http.StatusConflict, dtos.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, "Notification sent")
}

func NewNotificationController(client clients.NotificationsAPiClient) NotificationController {
	return &notificationController{client: client}
}
