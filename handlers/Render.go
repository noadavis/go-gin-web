package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Render(ctx *gin.Context, data gin.H, template string) {
	ctx.HTML(http.StatusOK, template, data)
}

func RenderJSON(ctx *gin.Context, data gin.H) {
	ctx.JSON(http.StatusOK, data["json"])
}

func RenderSTRING(ctx *gin.Context, data gin.H, content_type string) {
	ctx.Data(http.StatusOK, content_type, []byte(data["string"].(string)))
}
