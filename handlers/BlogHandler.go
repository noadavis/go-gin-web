package handlers

import (
	"go-gin-web/models"
	"go-gin-web/models/blog"
	"strconv"

	"github.com/gin-gonic/gin"
)

// [/blog/]
func BlogMainPage(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "Blog"
	appData.PageIcon = "edit"

	categories := blog.GetCategories(appData.UserAuth)

	records := blog.GetAllRecords(appData.UserAuth)

	Render(ctx, gin.H{"AppData": appData, "Categories": categories, "Records": records}, "page-blog.html")
}

// [/blog/:category_id]
func BlogPage_Category(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "Blog"
	appData.PageIcon = "edit"

	categories := blog.GetCategories(appData.UserAuth)

	records := []blog.BlogRecord{}
	if categoryId, err := strconv.Atoi(ctx.Param("category_id")); err == nil {
		records = blog.GetCategory(categoryId, appData.UserAuth)
	}

	Render(ctx, gin.H{"AppData": appData, "Categories": categories, "Records": records}, "page-blog.html")
}

// [/blog/record/:record_id]
func BlogPage_Record(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "Blog"
	appData.PageIcon = "edit"

	categories := blog.GetCategories(appData.UserAuth)

	author := false
	record := blog.BlogRecord{}
	if recordId, err := strconv.Atoi(ctx.Param("record_id")); err == nil {
		record = blog.GetRecord(recordId, appData.UserAuth, 0)
		if record.Id > 0 {
			appData.PageTitle = record.Name
			author = record.Author == appData.UserData.Id
		}
	}

	Render(ctx, gin.H{"AppData": appData, "Categories": categories, "Record": record, "Author": author}, "page-blog-record.html")
}

// [/blog/record/:record_id/edit]
func BlogPage_RecordEdit(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "Blog"
	appData.PageIcon = "edit"

	categories := blog.GetCategories(appData.UserAuth)

	author := false
	record := blog.BlogRecord{}
	if recordId, err := strconv.Atoi(ctx.Param("record_id")); err == nil {
		record = blog.GetRecord(recordId, appData.UserAuth, appData.UserData.Id)
		if record.Id > 0 {
			appData.PageTitle = record.Name
			author = record.Author == appData.UserData.Id
		}
	}

	Render(ctx, gin.H{"AppData": appData, "Categories": categories, "Record": record, "Author": author}, "page-blog-record-edit.html")
}
