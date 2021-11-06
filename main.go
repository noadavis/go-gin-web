// main.go

package main

import (
	"go-gin-web/models"
	"html/template"
	"log"
	"os"
	"time"

	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var props models.ConfigManager
var router *gin.Engine

func main() {
	// test config.json
	if !props.CheckConfig() {
		os.Exit(0)
	}

	// init database connection
	models.InitDb()

	// init object with app base information
	models.InitAppData()

	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()

	// create cookie store for auth
	store := cookie.NewStore([]byte(props.GetProps().AppConf.Salt))
	router.Use(sessions.Sessions("go_gin_web", store))

	// register some functions for templates
	router.SetFuncMap(template.FuncMap{
		"dateFormat": dateFormat,
	})

	// load template files
	router.LoadHTMLGlob("templates/*/*.html")

	// init routes
	InitRoutePaths()

	// run application with port from config.json
	log.Printf("http://127.0.0.1:%d\n", props.GetProps().AppConf.Port)
	router.Run(fmt.Sprintf(":%d", props.GetProps().AppConf.Port))
}

func dateFormat(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%02d/%02d/%d", day, month, year)
}
