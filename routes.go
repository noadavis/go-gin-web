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
	// static files from "./media" folder
	router.Static("/media", "./media")

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
	userRoutes.POST("/avatar", handlers.ShowUserPage_Avatar)

	// handle system route
	systemRoutes := router.Group("/system")
	systemRoutes.GET("/users", handlers.SystemUsersPage)
	systemRoutes.GET("/menu", handlers.SystemMenuPage)
	systemRoutes.GET("/menu/:menu_id", handlers.SystemMenuPage_Edit)
	systemRoutes.POST("/users/save", handlers.SystemUsersPage_Save)
	systemRoutes.GET("/users/new", handlers.SystemUsersPage_New)
	systemRoutes.GET("/users/edit/:user_id", handlers.SystemUsersPage_Edit)
	systemRoutes.GET("/users/:action_id/:user_id/:block", handlers.SystemUsersPage_Action)

	// handle blog route
	blogRoutes := router.Group("/blog")
	blogRoutes.GET("/", handlers.BlogMainPage)
	blogRoutes.GET("/record/:record_id", handlers.BlogPage_Record)
	blogRoutes.GET("/record/:record_id/edit", handlers.BlogPage_RecordEdit)
	blogRoutes.GET("/:category_id", handlers.BlogPage_Category)
}
