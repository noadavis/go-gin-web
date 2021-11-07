// models.user.go

package user

import (
	"go-gin-web/models"
	"log"
	"strconv"
)

type UserPermission struct {
	Id        int `db:"id"`
	User      int `db:"user"`
	Id_user   int `db:"id_user"`
	Id_editor int `db:"id_editor"`
	Id_admin  int `db:"id_admin"`
}

type UserListPermission struct {
	Id        int    `db:"id"`
	Fullname  string `db:"fullname"`
	Id_user   int    `db:"id_user"`
	Id_editor int    `db:"id_editor"`
	Id_admin  int    `db:"id_admin"`
	Enabled   int    `db:"enabled"`
}

func GetUserList() []UserListPermission {
	q := models.GetConnection()
	userList := []UserListPermission{}
	if err := q.Select(&userList, `SELECT us.id, us.fullname, up.id_user, up.id_editor, up.id_admin, lo.enabled
		FROM user_permission AS up
		LEFT JOIN user AS us ON up.user = us.id
		LEFT JOIN user_login AS lo ON lo.user = us.id
		ORDER BY us.fullname`); err != nil {
		log.Printf("GetUserList: %s\n", err.Error())
	}
	return userList
}

func GetUserInfo(userId int) models.UserInfo {
	q := models.GetConnection()
	userInfo := models.UserInfo{}
	// find user with username, password hash and not locked
	if err := q.Get(&userInfo, `SELECT
		us.id, us.login, us.fullname, us.email, up.id_admin, up.id_editor, up.id_user, lo.enabled
		FROM user AS us
		LEFT JOIN user_permission AS up ON up.user = us.id
		LEFT JOIN user_login AS lo ON lo.user = us.id
		WHERE us.id = ?`, userId); err != nil {
		log.Printf("GetUserInfo: %s\n", err.Error())
	}
	return userInfo
}

func LoginExists(login string) bool {
	q := models.GetConnection()
	userId := 0
	if err := q.Get(&userId, "SELECT id FROM user WHERE login = ?", login); err != nil {
		log.Printf("LoginExists [%s]: %s\n", login, err.Error())
		return false
	}
	return userId != 0
}

func RegNewUser(userInfo models.UserForm) bool {
	enabled, _ := strconv.Atoi(userInfo.Enabled)
	q := models.GetConnection()
	var err error
	userId := 0
	if _, err := q.NamedExec(`INSERT INTO user (login, fullname, email, avatar)
		VALUES (:ul, :un, :ue, 'avatar.png')`,
		map[string]interface{}{
			"ul": userInfo.Login,
			"un": userInfo.Fullname,
			"ue": userInfo.Email}); err != nil {
		log.Printf("RegNewUser [user]: %s\n", err.Error())
		return false
	}
	if err = q.Get(&userId, "SELECT id FROM user WHERE login = ?", userInfo.Login); err != nil {
		log.Printf("RegNewUser [user]: %s\n", err.Error())
		return false
	}
	if _, err = q.NamedExec(`INSERT INTO user_login (user, password, session, session_date, enabled)
		VALUES (:us, :ps, 0, 0, :ue)`,
		map[string]interface{}{
			"us": userId,
			"ps": calcPassword(userInfo.Password),
			"ue": enabled}); err != nil {
		log.Printf("RegNewUser [user_login]: %s\n", err.Error())
		return false
	}
	if _, err := q.NamedExec(`INSERT INTO user_permission (user, id_user, id_editor, id_admin)
		VALUES (:us, :iu, :ie, :ia)`,
		map[string]interface{}{
			"us": userId,
			"iu": findPermission("id_user", userInfo.Permissions),
			"ie": findPermission("id_editor", userInfo.Permissions),
			"ia": findPermission("id_admin", userInfo.Permissions)}); err != nil {
		log.Printf("RegNewUser [user_permission]: %s\n", err.Error())
		return false
	}
	return true
}

func EditUser(userInfo models.UserForm) bool {
	enabled, _ := strconv.Atoi(userInfo.Enabled)
	q := models.GetConnection()
	var err error
	userId := 0
	if err = q.Get(&userId, "SELECT id FROM user WHERE id = ?", userInfo.User); err != nil {
		log.Printf("EditUser: %s\n", err.Error())
		return false
	}
	if userId == 0 {
		return false
	}
	if _, err = q.NamedExec(`UPDATE user SET login = :ul, fullname = :un, email = :ue
		WHERE id = :us`,
		map[string]interface{}{
			"ul": userInfo.Login,
			"un": userInfo.Fullname,
			"ue": userInfo.Email,
			"us": userId}); err != nil {
		log.Printf("EditUser [user]: %s\n", err.Error())
		return false
	}
	if userInfo.Password == "" {
		_, err = q.NamedExec(`UPDATE user_login SET enabled = :ue
			WHERE id = :us`,
			map[string]interface{}{
				"ue": enabled,
				"us": userId})
	} else {
		_, err = q.NamedExec(`UPDATE user_login SET password = :ps, enabled = :ue
			WHERE id = :us`,
			map[string]interface{}{
				"ps": calcPassword(userInfo.Password),
				"ue": enabled,
				"us": userId})
	}
	if err != nil {
		log.Printf("EditUser [user_login]: %s\n", err.Error())
		return false
	}
	if _, err := q.NamedExec(`UPDATE user_permission SET id_user = :iu, id_editor = :ie, id_admin = :ia
		WHERE id = :us`,
		map[string]interface{}{
			"iu": findPermission("id_user", userInfo.Permissions),
			"ie": findPermission("id_editor", userInfo.Permissions),
			"ia": findPermission("id_admin", userInfo.Permissions),
			"us": userId}); err != nil {
		log.Printf("EditUser [user_permission]: %s\n", err.Error())
		return false
	}
	return true
}

func findPermission(permission string, permissions []string) string {
	for _, element := range permissions {
		if element == permission {
			return "1"
		}
	}
	return "0"
}

func BlockUser(block bool, userId int) bool {
	q := models.GetConnection()
	status := 1
	if block {
		status = 0
	}
	if _, err := q.NamedExec(`UPDATE user_login SET enabled = :status WHERE user = :user`,
		map[string]interface{}{
			"status": status,
			"user":   userId}); err != nil {
		log.Printf("BlockUser: %s\n", err.Error())
		return false
	}
	return true
}

func DeleteUser(userId int) bool {
	q := models.GetConnection()
	user := map[string]interface{}{"user": userId}
	var err error
	if _, err := q.NamedExec(`DELETE FROM user_permission WHERE user = :user`, user); err != nil {
		log.Printf("DeleteUser: %s\n", err.Error())
		return false
	}
	if _, err = q.NamedExec(`DELETE FROM user_login WHERE user = :user`, user); err != nil {
		log.Printf("DeleteUser: %s\n", err.Error())
		return false
	}
	if _, err = q.NamedExec(`DELETE FROM user WHERE id = :user`, user); err != nil {
		log.Printf("DeleteUser: %s\n", err.Error())
		return false
	}
	return true
}

func EditUserInfo(action, value string, userId int) bool {
	q := models.GetConnection()
	switch action {
	case "fullname":
		if _, err := q.NamedExec(`UPDATE user SET fullname = :fullname WHERE id = :user`,
			map[string]interface{}{
				"fullname": value,
				"user":     userId}); err != nil {
			log.Printf("EditUserInfo: %s\n", err.Error())
			return false
		}
	case "email":
		if _, err := q.NamedExec(`UPDATE user SET email = :email WHERE id = :user`,
			map[string]interface{}{
				"email": value,
				"user":  userId}); err != nil {
			log.Printf("EditUserInfo: %s\n", err.Error())
			return false
		}
	case "password":
		if _, err := q.NamedExec(`UPDATE user_login SET password = :password WHERE user = :user`,
			map[string]interface{}{
				"password": calcPassword(value),
				"user":     userId}); err != nil {
			log.Printf("EditUserInfo: %s\n", err.Error())
			return false
		}
	default:
		return false
	}
	return true
}

func SetNewAvatar(avatar string, userId int) bool {
	q := models.GetConnection()
	if _, err := q.NamedExec(`UPDATE user SET avatar = :avatar WHERE id = :user`,
		map[string]interface{}{
			"avatar": avatar,
			"user":   userId}); err != nil {
		log.Printf("SetNewAvatar: %s\n", err.Error())
		return false
	}
	return true
}
