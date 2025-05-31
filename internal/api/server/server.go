package server

import (
	"accounts/internal/api/health"
	"accounts/internal/api/router"
	"accounts/internal/api/v1/emails"
	"accounts/internal/api/v1/oauth_logins"
	refreshtokens "accounts/internal/api/v1/refresh_tokens"
	"accounts/internal/api/v1/roles"
	"accounts/internal/api/v1/users"
	"os"

	"accounts/internal/common/middlewares"
	"accounts/internal/core/settings"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

var Server *fiber.App

func Run() {

	app := setUpRouter()

	if _, inLambda := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); inLambda {
		fmt.Println("Running in Lambda")
		lambda.Start(ginadapter.NewV2(app).ProxyWithContext)
		return
	}

	app.Run(fmt.Sprintf(":%d", settings.Settings.PORT))
}

func setUpRouter() *gin.Engine {

	app := router.NewRouter()

	app.Use(middlewares.TraceMiddleware())
	//app.Use(middlewares.CatcherMiddleware)
	app.Use(middlewares.LoggerMiddleware())

	health.SetupHealthModule(app)
	roles.SetupRolesModule(app)
	users.SetupUsersModule(app)
	emails.SetupEmailsModule(app)
	refreshtokens.SetupRefreshTokensModule(app)
	oauth_logins.SetupOAuthModule(app)
	return app
}
