package controllers

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TestController interface {
	GetTest(ctx *gin.Context)
	PostTest(ctx *gin.Context)
}

type testController struct {
	testService services.TestService
}

func NewTestController(testService services.TestService) TestController {
	return &testController{testService}
}

// GetTest godoc
//
//	@Summary	Get a greeting
//	@Schemes
//	@Description	Get a greeting
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string
//	@Router			/test [get]
func (c *testController) GetTest(ctx *gin.Context) {
	greeting := c.testService.GetTest()
	ctx.JSON(http.StatusOK, gin.H{"greeting": greeting})
}

// PostTest godoc
//
//	@Summary		Post a greeting
//	@Schemes		dtos.GreetingTestDto
//	@Description	Post a greeting
//	@Accept			json
//	@Success		204
//	@Failure		400	{object}	dtos.ErrorResponse
//	@Router			/test [post]
func (c *testController) PostTest(ctx *gin.Context) {
	var greeting dtos.GreetingTestDto
	err := ctx.ShouldBindJSON(&greeting)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(err))
		return
	}
	c.testService.PostTest(greeting)
	ctx.JSON(http.StatusNoContent, nil)
}
