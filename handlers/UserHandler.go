package handlers

import (
	"go-gin-web/models"
	"go-gin-web/models/user"
	"time"

	"github.com/gin-gonic/gin"
)

// [/user/info]
func ShowUserInfoPage(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "User Info Page"
	appData.PageIcon = "info"

	userInfo := user.GetUserInfo(appData.UserData.Id)

	Render(ctx, gin.H{"AppData": appData, "Data": userInfo, "DateNow": time.Now()}, "page-user-info.html")
}

// [/user/user]
func ShowUserPage(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "User Page"
	appData.PageIcon = "lock"
	template, _ := models.CheckPermission(appData.Permissions, "id_user", "page-user-user.html")

	Render(ctx, gin.H{"AppData": appData}, template)
}

// [/user/editor]
func ShowEditorPage(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "Editor Page"
	appData.PageIcon = "lock"
	template, _ := models.CheckPermission(appData.Permissions, "id_editor", "page-user-editor.html")

	Render(ctx, gin.H{"AppData": appData}, template)
}
