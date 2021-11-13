package models

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type AppData struct {
	Menu        []MenuSection
	FullName    string
	ShortName   string
	Footer      FooterData
	BackPage    string
	UserAuth    bool
	Permissions []string
	PageTitle   string
	PageIcon    string
	UserData    UserData
}
type FooterData struct {
	Right string
	Left  string
}

var appData = AppData{}
var Salt = ""
var CookieLifetime = 0

func InitAppData() {
	appData.FullName = props.GetProps().AppConf.FullName
	appData.ShortName = props.GetProps().AppConf.ShortName
	appData.Footer.Right = "Right"
	appData.BackPage = "/"
	appData.UserAuth = false

	Salt = props.GetProps().AppConf.Salt
	CookieLifetime = props.GetProps().AppConf.CookieLifetime
}

func GetAppData() AppData {
	UpdateAppData()
	return appData
}

func UpdateAppData() {
	appData.Footer.Left = fmt.Sprintf("%d Left", time.Now().Year())
}

func GetSalt() string {
	return Salt
}

func GetCookieLifetime() int {
	return CookieLifetime
}

func GetVariableType(val interface{}) string {
	return fmt.Sprintf("%s", reflect.TypeOf(val))
}

func CheckPermission(permissions []string, needed string, template string) (string, bool) {
	for _, permission := range permissions {
		if permission == "id_admin" {
			return template, false
		} else {
			if permission == needed {
				return template, false
			}
		}
	}
	return "permission.html", true
}

func GetModifyMessage(error bool, messages []string) ModifyFormData {
	data := ModifyFormData{}
	data.Messages = messages
	if error {
		data.Title = "Runtime error"
		data.Action = "danger"
	} else {
		data.Title = "Data updated successfully"
		data.Action = "success"
	}
	return data
}

func CheckInt(value string) bool {
	if _, err := strconv.Atoi(value); err == nil {
		return true
	}
	return false
}
