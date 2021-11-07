package handlers

import (
	"fmt"
	"go-gin-web/models"
	"go-gin-web/models/user"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
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

// [/user/edituserinfo]
func ShowUserPage_EditInfo(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)

	answer := struct {
		Error bool   `json:"error"`
		Desc  string `json:"desc"`
	}{
		Error: true,
		Desc:  "",
	}

	action := ctx.DefaultPostForm("action", "")
	value := ctx.DefaultPostForm("value", "")

	if action == "" || value == "" {
		answer.Desc = "form error"
	} else {
		if answer.Error = !user.EditUserInfo(action, value, appData.UserData.Id); answer.Error {
			answer.Desc = "database error"
		}
	}

	RenderJSON(ctx, gin.H{"json": answer})
}

// [/user/avatar]
func ShowUserPage_Avatar(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)

	answer := struct {
		Error  bool   `json:"error"`
		Desc   string `json:"desc"`
		Avatar string `json:"avatar"`
	}{
		Error:  true,
		Desc:   "",
		Avatar: "",
	}

	if avatar, err := ctx.FormFile("avatar"); err == nil {
		avatarExt := filepath.Ext(avatar.Filename)
		avatarId := ksuid.New()
		answer.Avatar = fmt.Sprintf("%s%s", avatarId.String(), avatarExt)
		fmt.Println(answer.Avatar)
		if err = ctx.SaveUploadedFile(avatar, "./media/avatars/"+answer.Avatar); err == nil {
			if answer.Error = !user.SetNewAvatar(answer.Avatar, appData.UserData.Id); answer.Error {
				answer.Desc = "database error"
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}

	RenderJSON(ctx, gin.H{"json": answer})
}
