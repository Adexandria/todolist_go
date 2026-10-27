package main

import (
	"SM/handlers"
	"SM/models"
	"SM/repositories/Utilities"
	"SM/services/Middlewares"

	"github.com/gin-gonic/gin"
)

func registerRoutes(router *gin.Engine, tasks *handlers.Handler,
	tokenManager *Utilities.TokenManager, users *handlers.UserHandler) {

	authn := Middlewares.AuthenticateMiddleware(tokenManager)
	authz := Middlewares.AuthorizeMiddleware()

	api := router.Group("/api")
	anonymous := api.Group("", Middlewares.AllowAnonymous(), authn, authz)
	guestOnly := api.Group("", Middlewares.DisallowAuthentication(), Middlewares.AllowAnonymous(), authn, authz)
	user := api.Group("", Middlewares.AttachRoles(models.Roles[0]), authn, authz)
	admin := api.Group("", Middlewares.AttachRoles(models.Roles[1]), authn, authz)

	// tasks
	anonymous.GET("/tasks", tasks.GetAllTasks)
	anonymous.POST("/tasks", tasks.CreateTask)
	anonymous.GET("/tasks/filter", tasks.FilterTasks)
	anonymous.GET("/tasks/search", tasks.SearchByTask)
	anonymous.GET("/tasks/:id", tasks.GetTaskById)
	anonymous.PUT("/tasks/:id", tasks.UpdateTask)
	anonymous.DELETE("/tasks/:id", tasks.DeleteTask)

	// auth
	guestOnly.POST("/auth/signup", users.SignUp)
	guestOnly.POST("/auth/login", users.Login)
	guestOnly.POST("/auth/logout", users.Logout)
	guestOnly.GET("/auth/confirm-email", users.ConfirmEmailToken)
	guestOnly.POST("/auth/access-token", users.GenerateAccessToken)
	guestOnly.POST("/auth/reset-password", users.ResetPassword)
	guestOnly.POST("/auth/verify-password", users.ResetPasswordByToken)
	guestOnly.POST("/auth/otp/verify", users.VerifyOTP)

	user.POST("/auth/2fa/setup", users.SetUpTwoFactorAuthentication)
	user.POST("/auth/2fa/disable", users.DisableTwoFactorAuthentication)

	user.POST("/auth/2fa/request-otp", users.RequestOTP)

	// account
	user.GET("/account", users.GetUserProfile)
	user.DELETE("/account", users.DeleteUserAccount)
	user.PUT("/account/password", users.ChangePassword)
	user.PUT("/account/email", users.ChangeEmail)
	user.PUT("/account/phonenumber", users.ChangePhonenumber)
	user.GET("/account/verify-phonenumber", users.VerifyPhonenumberToken)
	user.PUT("/account/firstname", users.UpdateFirstName)
	user.PUT("/account/lastname", users.UpdateLastName)
	user.PUT("/account/username", users.UpdateUsername)

	// admin
	admin.GET("/admin/users/:id", users.GetUserById)
	admin.DELETE("/admin/users/:id", users.DeleteUserById)
	admin.PUT("/admin/users/:id/disable-lockout", users.DisableUserLockout)
	admin.POST("/admin/create", users.CreateAdmin)
	admin.POST("/admin/roles", users.ChangeRole)
	admin.POST("/admin/users/:id/disable-2fa", users.DisableUserTwoFactorAuthentication)
}
