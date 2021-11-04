// models.article.go

package menu

import (
	"fmt"
	"go-gin-web/models"
	"strings"
)

func checkAdmin(permissions []string) bool {
	for _, permission := range permissions {
		if permission == "id_admin" {
			return true
		}
	}
	return false
}

func getSqlPermissions(permissions []string, auth bool) string {
	if checkAdmin(permissions) {
		return ""
	}
	permissions = append(permissions, "")
	if auth {
		permissions = append(permissions, "id_auth")
	}
	return fmt.Sprintf(" AND permission in ('%s') ", strings.Join(permissions, "','"))
}

func GetMenu(permissions []string, auth, enabled bool, menuId int) []models.MenuSection {
	q := models.GetConnection()
	sections := []models.MenuSection{}
	showEnabled := "AND enabled = 1"
	if !enabled {
		showEnabled = ""
	}
	getMenuById := ""
	if menuId > 0 {
		getMenuById = fmt.Sprintf("AND id = %d", menuId)
	}
	permissoinsJoin := getSqlPermissions(permissions, auth)
	err := q.Select(&sections, fmt.Sprintf(`SELECT id, name, url, icon, permission, menu, enabled, ordering
		FROM menu WHERE parent = 0 %s %s %s
		ORDER BY ordering`, showEnabled, permissoinsJoin, getMenuById))
	if err != nil {
		fmt.Printf("GetMenu getSections: %s", err.Error())
	} else if len(sections) > 0 {
		for index, section := range sections {
			if section.MenuType == "multi" {
				childs := []models.MenuChild{}
				err = q.Select(&childs, fmt.Sprintf(`SELECT name, url, icon, permission, enabled, ordering
					FROM menu WHERE parent > 0 AND parent = %d %s %s
					ORDER BY ordering`, section.Id, showEnabled, permissoinsJoin))
				if err != nil {
					fmt.Printf("GetMenu getChilds: %s", err.Error())
				} else {
					if len(childs) > 0 {
						sections[index].Childs = childs
					}
				}
			}
		}
	}
	return sections
}
