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

	Render(ctx, gin.H{"AppData": appData, "Category": 0, "Categories": categories, "Records": records}, "page-blog.html")
}

func getCategoryName(categoryId int, categories []blog.BlogCategory) string {
	for _, element := range categories {
		if element.Id == categoryId {
			return element.Name
		}
	}
	return "Blog"
}

// [/blog/:category_id]
func BlogPage_Category(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "Blog"
	appData.PageIcon = "edit"

	categories := blog.GetCategories(appData.UserAuth)

	records := []blog.BlogRecord{}
	category := 0
	if categoryId, err := strconv.Atoi(ctx.Param("category_id")); err == nil {
		records = blog.GetCategoryContent(categoryId, appData.UserAuth)
		appData.PageTitle = getCategoryName(categoryId, categories)
		category = categoryId
	}

	Render(ctx, gin.H{"AppData": appData, "Category": category, "Categories": categories, "Records": records}, "page-blog.html")
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

	Render(ctx, gin.H{"AppData": appData, "Category": 0, "Categories": categories, "Record": record, "Author": author}, "page-blog-record.html")
}

// [/blog/record/:record_id/edit]
func BlogPage_RecordEdit(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "Blog"
	appData.PageIcon = "edit"

	categories := []blog.BlogCategory{}
	author := false
	record := blog.BlogRecord{}

	template, err := models.CheckPermission(appData.Permissions, "id_editor", "page-blog-record-edit.html")
	if !err {
		categories = blog.GetCategories(appData.UserAuth)
		if recordId, err := strconv.Atoi(ctx.Param("record_id")); err == nil {
			record = blog.GetRecord(recordId, appData.UserAuth, appData.UserData.Id)
			if record.Id > 0 {
				appData.PageTitle = record.Name
				author = record.Author == appData.UserData.Id
			}
		}
	}

	Render(ctx, gin.H{"AppData": appData, "Categories": categories, "Record": record, "Author": author}, template)
}

// [/blog/record/new]
func BlogPage_RecordNew(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "Blog: new record"
	appData.PageIcon = "edit"

	categories := []blog.BlogCategory{}
	record := blog.BlogRecord{}

	template, err := models.CheckPermission(appData.Permissions, "id_editor", "page-blog-record-edit.html")
	if !err {
		categories = blog.GetCategories(appData.UserAuth)
		// for create new record id: -1
		record.Id = -1
	}

	Render(ctx, gin.H{"AppData": appData, "Categories": categories, "Record": record, "Author": true}, template)
}

// [/blog/record/save]
func BlogPage_RecordSave(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)

	formData := blog.FormRecord{}
	ctx.ShouldBind(&formData)
	jsonAnswer := `{ "error": true, "desc": "" }`
	// check permissions: create, edit and delete can only id_editor and id_admin
	if _, err := models.CheckPermission(appData.Permissions, "id_editor", ""); !err {
		if models.CheckInt(formData.Id) {
			if formData.Action == "delete" {
				if blog.DeleteRecord(appData.UserData.Id, formData) {
					jsonAnswer = `{ "error": false, "desc": "" }`
				}
			} else {
				if blog.SaveRecord(appData.UserData.Id, formData) {
					jsonAnswer = `{ "error": false, "desc": "" }`
				}
			}
		}
	}

	RenderSTRING(ctx, gin.H{"string": jsonAnswer}, "application/json; charset=utf-8")
}

// [/blog/category/:category_id/edit]
func BlogPage_CategoryEdit(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "Blog"
	appData.PageIcon = "edit"

	category := blog.BlogCategory{}
	template, _ := models.CheckPermission(appData.Permissions, "id_editor", "page-blog-category-edit.html")
	if categoryId, err := strconv.Atoi(ctx.Param("category_id")); err == nil {
		category = blog.GetCategory(categoryId)
	}

	Render(ctx, gin.H{"AppData": appData, "Category": category}, template)
}

// [/blog/category/new]
func BlogPage_CategoryNew(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)
	appData.PageTitle = "Blog: new category"
	appData.PageIcon = "edit"

	category := blog.BlogCategory{}
	category.Id = -1
	template, _ := models.CheckPermission(appData.Permissions, "id_editor", "page-blog-category-edit.html")

	Render(ctx, gin.H{"AppData": appData, "Category": category}, template)
}

// [/blog/category/save]
func BlogPage_CategorySave(ctx *gin.Context) {
	appDataInterface, _ := ctx.Get("AppData")
	appData := appDataInterface.(models.AppData)

	formData := blog.FormCategory{}
	ctx.ShouldBind(&formData)
	jsonAnswer := `{ "error": true, "desc": "" }`
	if _, err := models.CheckPermission(appData.Permissions, "id_editor", ""); !err {
		if models.CheckInt(formData.Id) {
			if formData.Action == "delete" {
				if blog.DeleteCategory(formData) {
					jsonAnswer = `{ "error": false, "desc": "" }`
				}
			} else {
				if blog.SaveCategory(formData) {
					jsonAnswer = `{ "error": false, "desc": "" }`
				}
			}
		}
	}

	RenderSTRING(ctx, gin.H{"string": jsonAnswer}, "application/json; charset=utf-8")
}
