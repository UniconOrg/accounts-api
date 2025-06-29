package controllers

import (
	"accounts/internal/api/v1/role_policies/interface/dtos"
	"accounts/internal/common/logger"
	"accounts/internal/common/requests"

	"github.com/gin-gonic/gin"
)

func (c *RolePoliciesController) Create(ctx *gin.Context) {
	entry := logger.FromContext(ctx)

	entry.Info("Creating role policies")

	dto := requests.GetDTO[dtos.CreateRolePoliciesDTO](ctx)
	if dto == nil {
		entry.Error("Invalid request")
		return
	}

	command := dto.ToCommand()

	role_policies := c.service.Create(ctx.Request.Context(), command)

	ctx.JSON(role_policies.StatusCode, role_policies.ToMap())
}
