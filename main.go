package main

import (
	"SM/handlers"
	"SM/models"
	"SM/repositories/Utilities"
	"SM/services"
	"SM/services/Middlewares"
	"log/slog"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AttachRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(Middlewares.RolesKey, roles)
		c.Next()
	}
}
func AllowAnonymous() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(Middlewares.IsAnonymousRouteKey, true)
		c.Next()
	}
}

func DisallowAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(Middlewares.DisallowAuthenticationKey, true)
	}
}

func main() {
	container := services.CreateContainer()

	err := container.Provide(handlers.TaskHandler)
	if err != nil {
		panic(err)
	}

	passwordRule := Utilities.AccessRule{}.PasswordRule.SetMaximumPasswordLength(7).
		SetMinimumPasswordLength(3).EnableSmallLetter().EnableCapitalLetter().EnableNumber().EnableSpecialCharacter()

	accessRule := Utilities.AccessRule{}.EnableEmailConfirmation().EnablePasswordValidation(passwordRule)

	err = container.Provide(func() *Utilities.AccessRule {
		return &accessRule
	})

	if err != nil {
		panic(err)
	}

	err = container.Provide(handlers.UserHandlerCon)
	if err != nil {
		panic(err)
	}

	err = container.Provide(Utilities.ValidatorCon)
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.StaticFile("/docs/openapi.yml", "./docs/openapi.yml")

	url := ginSwagger.URL("/docs/openapi.yml")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	err = container.Invoke(func(h *handlers.Handler, jsonHandler *slog.JSONHandler,
		tokenManager *Utilities.TokenManager, a *handlers.UserHandler) {
		logger := slog.New(jsonHandler)

		logger.Info("Application is running")

		publicRoutes := router.Group("/api/tasks")

		publicRoutes.Use(Middlewares.AuthenticateMiddleware(*tokenManager),
			Middlewares.AuthorizeMiddleware())

		publicRoutes.GET("/:id", AllowAnonymous(), h.GetTaskById)
		publicRoutes.POST("/", AllowAnonymous(), h.CreateTask)
		publicRoutes.PUT("/:id", AllowAnonymous(), h.UpdateTask)
		publicRoutes.DELETE("/:id", AllowAnonymous(), h.DeleteTask)
		publicRoutes.GET("/", AllowAnonymous(), h.GetAllTasks)
		publicRoutes.GET("/filter", AllowAnonymous(), h.FilterTasks)
		publicRoutes.GET("/search", AllowAnonymous(), h.SearchByTask)

		authRoutes := router.Group("/api/auth")

		authRoutes.Use(Middlewares.AuthenticateMiddleware(*tokenManager),
			Middlewares.AuthorizeMiddleware())

		authRoutes.POST("/signup", AttachRoles(models.Roles[0]), a.SignUp)
		authRoutes.POST("/login", AllowAnonymous(), DisallowAuthentication(), a.Login)
		authRoutes.POST("/logout", AllowAnonymous(), DisallowAuthentication(), a.Logout)
		authRoutes.POST("/2fa/setup", AttachRoles(models.Roles[0]), a.SetUpTwoFactorAuthentication)
		authRoutes.POST("/2fa/disable", AttachRoles(models.Roles[0]), a.DisableTwoFactorAuthentication)
		authRoutes.GET("/otp/verify", AttachRoles(models.Roles[0]), a.VerifyOTP)
		authRoutes.POST("/reset-password", AttachRoles(models.Roles[0]), a.ResetPassword)
		authRoutes.POST("/reset-password-token", AttachRoles(models.Roles[0]), a.ResetPasswordByToken)
		authRoutes.GET("/confirm-email", AttachRoles(models.Roles[0]), a.ConfirmEmailToken)

		accountRoutes := router.Group("/api/account")
		accountRoutes.Use(Middlewares.AuthenticateMiddleware(*tokenManager),
			Middlewares.AuthorizeMiddleware())

		accountRoutes.PUT("/password", AttachRoles(models.Roles[0]), a.ChangePassword)

		accountRoutes.PUT("/email", AttachRoles(models.Roles[0]), a.ChangeEmail)

		accountRoutes.PUT("/phonenumber", AttachRoles(models.Roles[0]), a.ChangePassword)
		accountRoutes.GET("/verify-phonenumber", AttachRoles(models.Roles[0]), a.ChangePhonenumber)
		accountRoutes.GET("/", AttachRoles(models.Roles[0]), a.GetUserProfile)
		accountRoutes.DELETE("/", AttachRoles(models.Roles[0]), a.DeleteUserAccount)
		accountRoutes.PUT("/firstname", AttachRoles(models.Roles[0]), a.UpdateFirstName)
		accountRoutes.PUT("/lastname", AttachRoles(models.Roles[0]), a.UpdateLastName)
		accountRoutes.PUT("/username", AttachRoles(models.Roles[0]), a.UpdateUsername)

		adminRoutes := router.Group("/api/admin")
		adminRoutes.Use(Middlewares.AuthenticateMiddleware(*tokenManager),
			Middlewares.AuthorizeMiddleware())

		adminRoutes.DELETE("/users/:id", AttachRoles(models.Roles[1]), a.DeleteUserById)
		adminRoutes.GET("/users/:id", AttachRoles(models.Roles[1]), a.GetUserById)
		adminRoutes.PUT("/users/:id/disable-lockout", AttachRoles(models.Roles[1]), a.DisableUserLockout)
		adminRoutes.POST("/create", AttachRoles(models.Roles[1]), a.CreateAdmin)
		adminRoutes.POST("/roles", AttachRoles(models.Roles[1]), a.ChangeRole)
	})

	if err != nil {
		panic(err)
	}

	err = router.Run(":8080")
	if err != nil {
		return
	}
}
