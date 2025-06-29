package controllers

import (
	"accounts/internal/common/logger"

	"github.com/gin-gonic/gin"
)

func (c *RolePoliciesController) Info(ctx *gin.Context) {
	entry := logger.FromContext(ctx)

	entry.Info("Getting role policies info")

	role_id := ctx.Param("role_id")

	role_policies := c.service.Info(ctx.Request.Context(), role_id)

	ctx.JSON(role_policies.StatusCode, role_policies.ToMap())
}
