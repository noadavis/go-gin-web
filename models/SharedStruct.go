package models

type MenuSection struct {
	Id         int    `db:"id"`
	Name       string `db:"name"`
	Url        string `db:"url"`
	Icon       string `db:"icon"`
	Permission string `db:"permission"`
	MenuType   string `db:"menu"`
	Childs     []MenuChild
	Enabled    int `db:"enabled"`
	Ordering   int `db:"ordering"`
}

type MenuChild struct {
	Name       string `db:"name"`
	Url        string `db:"url"`
	Icon       string `db:"icon"`
	Permission string `db:"permission"`
	Enabled    int    `db:"enabled"`
	Ordering   int    `db:"ordering"`
}

type SystemUser struct {
	Id       int    `db:"id"`
	Login    string `db:"login"`
	Fullname string `db:"fullname"`
	Avatar   string `db:"avatar"`
}

type UserInfo struct {
	Id        int    `db:"id"`
	Login     string `db:"login"`
	Fullname  string `db:"fullname"`
	Email     string `db:"email"`
	Id_admin  int    `db:"id_admin"`
	Id_editor int    `db:"id_editor"`
	Id_user   int    `db:"id_user"`
	Enabled   int    `db:"enabled"`
}

type ModifyFormData struct {
	Title    string
	Action   string
	Messages []string
}

type UserForm struct {
	Permissions []string `form:"permissions[]"`
	Enabled     string   `form:"enabled"`
	Login       string   `form:"login"`
	Fullname    string   `form:"fullname"`
	Email       string   `form:"email"`
	Password    string   `form:"password"`
	Passcheck   string   `form:"passcheck"`
	User        string   `form:"user"`
	Action      string   `form:"action"`
}
