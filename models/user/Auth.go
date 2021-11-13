package user

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go-gin-web/models"
	"strconv"
)

func calcPassword(pass string) string {
	passWithSalt := pass + models.GetSalt()
	hash := md5.Sum([]byte(passWithSalt))
	return hex.EncodeToString(hash[:])
}

func checkPassword(pass, passhash string) bool {
	return calcPassword(pass) == passhash
}

func CheckSession(sessionId, sessionName, until string) models.UserData {
	q := models.GetConnection()
	user := models.UserData{}
	if err := q.Get(&user, `SELECT us.id, us.login, us.fullname, us.avatar FROM user_login AS lo
		LEFT JOIN user AS us ON lo.user = us.id
		WHERE us.login = ? AND lo.session = ? AND lo.session_date > ? AND lo.enabled = 1`,
		sessionName, sessionId, until); err != nil {
		fmt.Printf("CheckSession: %s\n", err.Error())
	}
	user.Editor = false
	return user
}

func FindActiveUserByPassword(username, password string) (bool, int) {
	q := models.GetConnection()
	login := 0
	passhash := calcPassword(password)
	// find user with username, password hash and not locked
	if err := q.Get(&login, `SELECT us.id FROM user_login AS lo
		LEFT JOIN user AS us ON lo.user = us.id
		WHERE us.login = ? AND lo.password = ? AND lo.enabled = 1`, username, passhash); err != nil {
		fmt.Printf("FindActiveUserByPassword: %s\n", err.Error())
		return false, 0
	}
	return true, login
}

func UpdateDbSession(userId int, now, session string) bool {
	// fmt.Printf("UpdateDbSession: %d %s %s\n", userId, now, session)
	q := models.GetConnection()
	if _, err := q.NamedExec(`UPDATE user_login SET session_date = :now, session = :session
		WHERE user = :user`,
		map[string]interface{}{
			"now":     now,
			"session": session,
			"user":    strconv.Itoa(userId)}); err != nil {
		fmt.Printf("UpdateDbSession: %s\n", err.Error())
		return false
	}
	return true
}

func getUserPermissions(userId int) UserPermission {
	q := models.GetConnection()
	permission := UserPermission{}
	if err := q.Get(&permission, `SELECT * FROM user_permission WHERE id = ?`, userId); err != nil {
		fmt.Printf("getUserPermissions: %s\n", err.Error())
	}
	return permission
}

func GetUserPermissionsMap(userId int) map[string]bool {
	permission := getUserPermissions(userId)
	permissions := map[string]bool{}
	if permission.Id > 0 {
		permissions["id_user"] = permission.Id_user == 1
		permissions["id_editor"] = permission.Id_editor == 1
		permissions["id_admin"] = permission.Id_admin == 1
	}
	return permissions
}

func GetUserPermissionsList(userId int) []string {
	permission := getUserPermissions(userId)
	permissions := []string{}
	if permission.Id > 0 {
		if permission.Id_user == 1 {
			permissions = append(permissions, "id_user")
		}
		if permission.Id_editor == 1 {
			permissions = append(permissions, "id_editor")
		}
		if permission.Id_admin == 1 {
			permissions = append(permissions, "id_admin")
		}
	}
	return permissions
}
