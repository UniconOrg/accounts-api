package controllers

import (
	"accounts/internal/api/v1/emails/interface/dtos"
	"accounts/internal/common/logger"
	"accounts/internal/common/requests"

	"github.com/gin-gonic/gin"
)

func (c *EmailsController) SetPassword(ctx *gin.Context) {

	entry := logger.FromContext(ctx)

	dto := requests.GetDTO[dtos.SetPasswordDTO](ctx)

	if dto == nil {
		entry.Error("Error al parsear el JSON")
		return
	}

	command := dto.ToCommand()

	token := requests.GetToken(ctx)
	if token == nil {
		entry.Error("Failed to get token from request")
		return
	}

	response := c.userService.SetPassword(ctx.Request.Context(), command, token.Token)
	// Se almacena el objeto para que el middleware lo procese
	ctx.JSON(response.StatusCode, response.ToMap())
}
