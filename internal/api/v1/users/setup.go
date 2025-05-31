package users

import (
	"accounts/internal/api/v1/users/domain/services"
	"accounts/internal/api/v1/users/interface/controllers"
	"accounts/internal/core/settings"
	roles "accounts/internal/db/postgres/role"
	users "accounts/internal/db/postgres/users"

	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupUsersModule(app *gin.Engine) {

	db, err := gorm.Open(postgres.Open(settings.Settings.POSTGRES_DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println(db)

	service := services.NewUsersService(
		users.NewUserPostgresRepository(db),
		roles.NewRolePostgresRepository(db),
	)

	controller := controllers.NewUsersController(*service)

	// Rutas de users
	group := app.Group(settings.Settings.ROOT_PATH + "/api/v1/users")

	group.POST("", controller.Create)
	group.GET("", controller.List)

}
