package handlers

import (
	"go-gin-web/models"
	"go-gin-web/models/menu"
	"go-gin-web/models/user"
	"strconv"

	"github.com/gin-gonic/gin"
)

// [/system/users/save]
func SystemUsersPage_Save(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.BackPage = "/system/users"

	messages := []string{}
	err := true
	template, error := models.CheckPermission(appData.Permissions, "id_admin", "modify.html")
	if !error {
		formData := models.UserForm{}
		ctx.ShouldBind(&formData)
		switch formData.Action {
		case "new":
			if user.RegNewUser(formData) {
				messages = append(messages, "User created")
				err = false
			}
		case "edit":
			if user.EditUser(formData) {
				messages = append(messages, "User edited")
				err = false
			}
		default:
			messages = append(messages, "Form error")
		}
	}
	ModData := models.GetModifyMessage(err, messages)

	Render(ctx, gin.H{"AppData": appData, "Message": ModData}, template)
}

// [/system/users]
func SystemUsersPage(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "System: Users"
	appData.PageIcon = "users"
	template, error := models.CheckPermission(appData.Permissions, "id_admin", "page-system-users.html")

	userList := []user.UserListPermission{}
	if !error {
		userList = user.GetUserList()
	}

	Render(ctx, gin.H{"AppData": appData, "Data": userList}, template)
}

// [/system/menu]
func SystemMenuPage(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "System: Menu"
	appData.PageIcon = "list"
	template, error := models.CheckPermission(appData.Permissions, "id_admin", "page-system-menu.html")

	menuList := []models.MenuSection{}
	if !error {
		menuList = menu.GetMenu([]string{"id_admin"}, false, false, 0)
	}

	Render(ctx, gin.H{"AppData": appData, "Data": menuList}, template)
}

// [/system/menu/:menu_id]
func SystemMenuPage_Edit(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "System: Edit Menu"
	appData.PageIcon = "list"
	template, error := models.CheckPermission(appData.Permissions, "id_admin", "page-system-menu-edit.html")

	menuItem := models.MenuSection{}
	if !error {
		if menuId, err := strconv.Atoi(ctx.Param("menu_id")); err == nil {
			menuList := menu.GetMenu([]string{"id_admin"}, false, false, menuId)
			if len(menuList) > 0 {
				menuItem = menuList[0]
			}
		}
	}

	Render(ctx, gin.H{"AppData": appData, "Data": menuItem}, template)
}

// [/system/users/:user_id]
func SystemUsersPage_Edit(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "System: Edit User"
	appData.PageIcon = "user"
	template, error := models.CheckPermission(appData.Permissions, "id_admin", "page-system-user.html")

	userInfo := models.UserInfo{}
	if !error {
		if userId, err := strconv.Atoi(ctx.Param("user_id")); err == nil {
			userInfo = user.GetUserInfo(userId)
		}
	}
	type UserForm struct {
		Show bool
		Type string
	}
	userForm := UserForm{Show: userInfo.Id > 0, Type: "edit"}

	Render(ctx, gin.H{"AppData": appData, "Data": userInfo, "UserForm": userForm}, template)
}

// [/system/users/new]
func SystemUsersPage_New(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "System: New User"
	appData.PageIcon = "user"
	template, _ := models.CheckPermission(appData.Permissions, "id_admin", "page-system-user.html")

	userInfo := models.UserInfo{}
	type UserForm struct {
		Show bool
		Type string
	}
	userForm := UserForm{Show: true, Type: "new"}

	Render(ctx, gin.H{"AppData": appData, "Data": userInfo, "UserForm": userForm}, template)
}

func getGetParam(value string) int {
	if converted, err := strconv.Atoi(value); err == nil {
		return converted
	}
	return 0
}

// [/system/users/:action_id/:user_id]
func SystemUsersPage_Action(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)

	if _, error := models.CheckPermission(appData.Permissions, "id_admin", "pass"); error {
		RenderSTRING(ctx, gin.H{"string": `{ "error": true, "desc": "permission error" }`},
			"application/json; charset=utf-8")
		return
	}

	actionId := getGetParam(ctx.Param("action_id"))
	userId := getGetParam(ctx.Param("user_id"))
	block := getGetParam(ctx.Param("block"))
	if actionId == 0 || userId == 0 {
		RenderSTRING(ctx, gin.H{"string": `{ "error": true, "desc": "form data error" }`},
			"application/json; charset=utf-8")
	} else {
		jsonAnswer := `{ "error": true, "desc": "" }`
		switch actionId {
		case 1:
			// user info
			userInfo := user.GetUserInfo(userId)
			RenderJSON(ctx, gin.H{"json": userInfo})
			return
		case 2:
			// enable/block user
			if user.BlockUser(block == 1, userId) {
				jsonAnswer = `{ "error": false, "desc": "" }`
			}
		case 3:
			// delete user
			if user.DeleteUser(userId) {
				jsonAnswer = `{ "error": false, "desc": "" }`
			}
		}
		RenderSTRING(ctx, gin.H{"string": jsonAnswer}, "application/json; charset=utf-8")
	}
}
