package handlers

import (
	"go-gin-web/models"
	"go-gin-web/models/user"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// [/]
func ShowIndexPage(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "Home Page"
	appData.PageIcon = "home"
	Render(ctx, gin.H{"AppData": appData}, "page-home.html")
}

// [/permission]
func ShowPermissionPage(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "Permission Page"
	appData.PageIcon = "lock"
	Render(ctx, gin.H{"AppData": appData}, "permission.html")
}

// [/auth/login]
func ShowLoginPage(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	getParams := ctx.Request.URL.Query()
	message := ""
	if val, ok := getParams["a"]; ok {
		if len(val) > 0 {
			switch val[0] {
			case "fail":
				message = "User or Password is not correct"
			}
		}
	}
	Render(ctx, gin.H{"AppData": appDataInterface, "Message": message}, "login.html")
}

// [/auth/register]
func ShowRegisterPage(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	getParams := ctx.Request.URL.Query()
	message := ""
	if val, ok := getParams["a"]; ok {
		if len(val) > 0 {
			switch val[0] {
			case "success":
				message = "New user created"
			case "fail":
				message = "An error occurred during registration"
			case "pass":
				message = "Password mismatch"
			}
		}
	}
	appData := appDataInterface.(models.AppData)
	Render(ctx, gin.H{"AppData": appData, "Message": message}, "register.html")
}

// post request for login or register [/auth/auth]
func Auth(ctx *gin.Context) {
	formData := models.UserForm{}
	ctx.ShouldBind(&formData)

	location := url.URL{Path: "/"}
	switch formData.Action {
	case "login":
		// perform login
		if find, userId := user.FindActiveUserByPassword(formData.Login, formData.Password); find {
			if updateUserSession(userId, formData.Login, sessions.Default(ctx)) {
				// login complete, redirect ot home
				location = url.URL{Path: "/"}
			}
		} else {
			location = url.URL{Path: "/auth/login", RawQuery: "a=fail"}
		}
	case "reg":
		// create new user
		if formData.Password == formData.Passcheck {
			formData.Enabled = "1"
			formData.Permissions = append(formData.Permissions, "id_user")
			if user.RegNewUser(formData) {
				// register complete
				location = url.URL{Path: "/auth/register", RawQuery: "a=success"}
			} else {
				location = url.URL{Path: "/auth/register", RawQuery: "a=fail"}
			}
		} else {
			location = url.URL{Path: "/auth/register", RawQuery: "a=pass"}
		}
	}
	ctx.Redirect(302, location.RequestURI())
}

// logout request [/auth/logout]
func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("gin_session_id")
	session.Delete("gin_username")
	session.Save()
	location := url.URL{Path: "/"}
	ctx.Redirect(302, location.RequestURI())
}

func updateUserSession(userId int, username string, session sessions.Session) bool {
	now := getUnixTimestampUntil()
	sessionId := newSessionId(25)
	session.Set("gin_session_id", sessionId)
	session.Set("gin_username", username)
	session.Options(sessions.Options{Path: "/", MaxAge: models.GetCookieLifetime()})
	session.Save()
	user.UpdateDbSession(userId, now, sessionId)
	return true
}

func newSessionId(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	count := len(letterBytes)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(count)]
	}
	return string(b)
}

func getUnixTimestampUntil() string {
	until := time.Now().Unix() + int64(models.GetCookieLifetime())
	return strconv.FormatInt(until, 10)
}
