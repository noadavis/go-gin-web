// routes.go

package main

import (
	"go-gin-web/handlers"
	"go-gin-web/middleware"
)

func InitRoutePaths() {
	// add middleware
	router.Use(middleware.SessionAuthMiddleware())

	// static files: js, css, images in "./static" folder
	router.Static("/static", "./static")

	// handle index route
	router.GET("/", handlers.ShowIndexPage)
	router.GET("/permission", handlers.ShowPermissionPage)

	// handle auth route
	authRoutes := router.Group("/auth")
	authRoutes.GET("/login", handlers.ShowLoginPage)
	authRoutes.GET("/logout", handlers.Logout)
	authRoutes.GET("/register", handlers.ShowRegisterPage)
	authRoutes.POST("/auth", handlers.Auth)

	// handle user route
	userRoutes := router.Group("/user")
	userRoutes.GET("/info", handlers.ShowUserInfoPage)
	userRoutes.GET("/user", handlers.ShowUserPage)
	userRoutes.GET("/editor", handlers.ShowEditorPage)
	userRoutes.POST("/edituserinfo", handlers.ShowUserPage_EditInfo)

	// handle system route
	systemnRoutes := router.Group("/system")
	systemnRoutes.GET("/users", handlers.SystemUsersPage)
	systemnRoutes.GET("/menu", handlers.SystemMenuPage)
	systemnRoutes.GET("/menu/:menu_id", handlers.SystemMenuPage_Edit)
	systemnRoutes.POST("/users/save", handlers.SystemUsersPage_Save)
	systemnRoutes.GET("/users/new", handlers.SystemUsersPage_New)
	systemnRoutes.GET("/users/edit/:user_id", handlers.SystemUsersPage_Edit)
	systemnRoutes.GET("/users/:action_id/:user_id/:block", handlers.SystemUsersPage_Action)
}
