package role_policies

import (
	"accounts/internal/api/v1/role_policies/domain/services"
	"accounts/internal/api/v1/role_policies/interface/controllers"
	"accounts/internal/core/settings"
	policies_pg "accounts/internal/db/postgres/policies"
	roles_pg "accounts/internal/db/postgres/role"
	role_policies_pg "accounts/internal/db/postgres/role_policies"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRolePoliciesModule(router *gin.Engine, db *gorm.DB) {
	// clients

	// repositories
	role_policies_repository := role_policies_pg.NewRolePoliciesPostgresRepository(db)
	role_repository := roles_pg.NewRolePostgresRepository(db)
	policies_repository := policies_pg.NewPoliciesPostgresRepository(db)

	// services
	role_policies_service := services.NewRolePoliciesService(role_policies_repository, role_repository, policies_repository)

	// controllers
	role_policies_controller := controllers.NewRolePoliciesController(role_policies_service)

	// routes

	role_policies_route := router.Group(settings.Settings.ROOT_PATH + "/api/v1/role_policies")

	role_policies_route.POST("", role_policies_controller.Create)
	role_policies_route.GET("/:role_id", role_policies_controller.Info)
}
