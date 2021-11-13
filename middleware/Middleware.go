package middleware

import (
	"fmt"
	"go-gin-web/models"
	"go-gin-web/models/menu"
	"go-gin-web/models/user"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var props models.ConfigManager
var openPaths = map[string]bool{
	"/":                       true,
	"/blog/":                  true,
	"/blog/:category_id":      true,
	"/blog/record/:record_id": true,
}
var noAuthPaths = map[string]bool{
	"/auth/login":    true,
	"/auth/register": true,
	"/auth/auth":     true,
	"/permission":    true,
}
var noMiddlewarePaths = map[string]bool{
	"/static/*filepath": true,
	"/media/*filepath":  true,
	"/favicon.ico":      true,
}

func getUnixTimestampNow() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentPath := ctx.FullPath()
		if noMiddlewarePaths[currentPath] {
			return
		}
		fmt.Println("SessionAuthMiddleware: currentPath", currentPath)
		// with JSON request there is no need to obtain additional information, such as menu
		jsonRequest := ctx.Request.Header.Get("Accept") == "application/json"
		// get AppData with basic app variables, as new object
		// update AppData for current user: menu, permissions, auth, userinfo
		appData := models.GetAppData()
		if !noAuthPaths[currentPath] {
			session := sessions.Default(ctx)
			session.Options(sessions.Options{Path: "/"})
			sessionId := session.Get("gin_session_id")
			sessionName := session.Get("gin_username")
			if sessionId != nil && sessionName != nil {
				appData.UserData = user.CheckSession(sessionId.(string), sessionName.(string), getUnixTimestampNow())
				appData.UserAuth = appData.UserData.Login != ""
				if appData.UserAuth {
					appData.Permissions = user.GetUserPermissionsList(appData.UserData.Id)
					if _, err := models.CheckPermission(appData.Permissions, "id_editor", ""); !err {
						appData.UserData.Editor = true
					}
				}
			}
			// redirect to permission page for not authorized users
			if !openPaths[currentPath] {
				if !appData.UserAuth {
					location := url.URL{Path: "/permission"}
					ctx.Redirect(302, location.RequestURI())
					return
				}
			}
		}
		fmt.Println("SessionAuthMiddleware: isAuth", appData.UserAuth)
		if !jsonRequest {
			appData.Menu = menu.GetMenu(appData.Permissions, appData.UserAuth, true, 0)
		}
		ctx.Set("AppData", appData)
		ctx.Next()
	}
}
