package controllers

import (
	"accounts/internal/api/v1/policies/interface/dtos"
	"accounts/internal/common/logger"
	"accounts/internal/common/requests"

	"github.com/gin-gonic/gin"
)

func (c *PoliciesController) Create(ctx *gin.Context) {

	entry := logger.FromContext(ctx)

	entry.Info("Creating policy")

	dto := requests.GetDTO[dtos.CreatePolicyDTO](ctx)
	if dto == nil {
		entry.Error("Invalid request")
		return
	}

	command := dto.ToCommand()

	policy := c.policies_service.Create(ctx.Request.Context(), command)

	ctx.JSON(policy.StatusCode, policy.ToMap())
}
