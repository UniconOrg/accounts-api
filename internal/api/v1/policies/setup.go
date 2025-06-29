package policies

import (
	"accounts/internal/api/v1/policies/domain/services"
	"accounts/internal/api/v1/policies/interface/controllers"
	"accounts/internal/core/settings"
	policies_pg "accounts/internal/db/postgres/policies"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupPoliciesModule(r *gin.Engine, db *gorm.DB) {

	// clients

	// repositories
	policies_repository := policies_pg.NewPoliciesPostgresRepository(db)

	// services
	policies_service := services.NewPoliciesService(policies_repository)

	// controllers
	policies_controller := controllers.NewPoliciesController(policies_service)

	// routes

	policies_route := r.Group(settings.Settings.ROOT_PATH + "/api/v1/policies")

	policies_route.POST("", policies_controller.Create)
	policies_route.GET("", policies_controller.List)
}
