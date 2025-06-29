package controllers

import (
	"accounts/internal/common/responses"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

func (c *PoliciesController) List(ctx *gin.Context) {

	policies, err := c.policies_service.List(ctx.Request.Context())
	if err != nil {
		ctx.JSON(fiber.StatusBadRequest, responses.Response{
			Status: fiber.StatusBadRequest,
			Errors: []string{err.Error()},
		})
		return
	}

	customResponse := responses.Response{
		Status: fiber.StatusOK,
		Data:   policies,
	}

	// Se almacena el objeto para que el middleware lo procese
	ctx.JSON(fiber.StatusOK, customResponse)
}
