package blog

import (
	"fmt"
	"go-gin-web/models"
	"log"
	"time"
)

type BlogCategory struct {
	Id    int    `db:"id"`
	Name  string `db:"name"`
	Alias string `db:"alias"`
	Auth  int    `db:"auth"`
}

type BlogRecord struct {
	Id          int       `db:"id"`
	CategoryId  int       `db:"category"`
	Category    string    `db:"cname"`
	Name        string    `db:"name"`
	Author      int       `db:"author"`
	Fullname    string    `db:"fullname"`
	Preview     string    `db:"preview"`
	Text        string    `db:"text"`
	Datecreated time.Time `db:"datecreated"`
	Datechanged time.Time `db:"datechanged"`
}

type FormRecord struct {
	Id       string `form:"record"`
	Title    string `form:"title"`
	Category string `form:"category"`
	Preview  string `form:"preview"`
	Content  string `form:"content"`
	Action   string `form:"action"`
}

type FormCategory struct {
	Id     string `form:"category"`
	Name   string `form:"name"`
	Alias  string `form:"alias"`
	Auth   string `form:"auth"`
	Action string `form:"action"`
}

func GetCategories(auth bool) []BlogCategory {
	q := models.GetConnection()
	categories := []BlogCategory{}
	showAll := ""
	if !auth {
		showAll = "WHERE auth = 0"
	}
	if err := q.Select(&categories, fmt.Sprintf(`SELECT id, name
		FROM category %s ORDER BY id`, showAll)); err != nil {
		log.Printf("GetCategories: %s", err.Error())
	}
	return categories
}

func GetCategory(categoryId int) BlogCategory {
	q := models.GetConnection()
	category := BlogCategory{}
	if err := q.Get(&category, `SELECT id, name, alias, auth 
		FROM category WHERE id = ?`, categoryId); err != nil {
		log.Printf("GetCategory: %s", err.Error())
	}
	return category
}

func SaveCategory(data FormCategory) bool {
	q := models.GetConnection()
	if data.Id == "-1" {
		// add new category
		if _, err := q.NamedExec(`INSERT INTO category (name, auth, alias, permission)
			VALUES (:na, :au, :al, '')`,
			map[string]interface{}{
				"na": data.Name,
				"au": data.Auth,
				"al": data.Alias}); err != nil {
			log.Printf("SaveCategory: %s\n", err.Error())
			return false
		}
	} else {
		if _, err := q.NamedExec(`UPDATE category SET name = :na, auth = :au, alias = :al, permission = ''
			WHERE id = :id`,
			map[string]interface{}{
				"na": data.Name,
				"au": data.Auth,
				"al": data.Alias,
				"id": data.Id}); err != nil {
			log.Printf("SaveCategory: %s\n", err.Error())
			return false
		}
	}
	return true
}

func GetCategoryContent(category int, auth bool) []BlogRecord {
	q := models.GetConnection()
	records := []BlogRecord{}
	showAll := ""
	if !auth {
		showAll = " AND c.auth = 0"
	}
	sql := fmt.Sprintf(`SELECT b.*, u.fullname, c.name AS cname FROM blog AS b
		LEFT JOIN category AS c ON b.category = c.id
		LEFT JOIN user AS u ON b.author = u.id
		WHERE c.id = ? %s ORDER BY b.datecreated DESC`, showAll)
	if err := q.Select(&records, sql, category); err != nil {
		log.Printf("GetRecordsById: %s\n", err.Error())
	}
	return records
}

func GetAllRecords(auth bool) []BlogRecord {
	q := models.GetConnection()
	records := []BlogRecord{}
	showAll := ""
	if !auth {
		showAll = "WHERE c.auth = 0"
	}
	sql := fmt.Sprintf(`SELECT b.*, u.fullname, c.name AS cname FROM blog AS b
		LEFT JOIN category AS c ON b.category = c.id
		LEFT JOIN user AS u ON b.author = u.id
		%s ORDER BY b.datecreated DESC`, showAll)
	if err := q.Select(&records, sql); err != nil {
		log.Printf("GetAllRecords: %s\n", err.Error())
	}
	return records
}

func GetRecord(recordId int, auth bool, author int) BlogRecord {
	q := models.GetConnection()
	record := BlogRecord{}
	showAll := ""
	if !auth {
		showAll = "AND c.auth = 0"
	}
	checkAuthor := ""
	if author > 0 {
		checkAuthor = fmt.Sprintf("AND b.author = %d", author)
	}
	sql := fmt.Sprintf(`SELECT b.*, u.fullname, c.name AS cname FROM blog AS b
		LEFT JOIN category AS c ON b.category = c.id
		LEFT JOIN user AS u ON b.author = u.id
		WHERE b.id = ? %s %s ORDER BY b.datecreated DESC`, showAll, checkAuthor)
	if err := q.Get(&record, sql, recordId); err != nil {
		log.Printf("GetAllRecords: %s\n", err.Error())
	}
	return record
}

func SaveRecord(userId int, data FormRecord) bool {
	q := models.GetConnection()
	if data.Id == "-1" {
		// add new record
		if _, err := q.NamedExec(`INSERT INTO blog (category, name, author, preview, text, datecreated, datechanged) VALUES
			(:ca, :na, :au, :pr, :te, :dcr, :dch)`,
			map[string]interface{}{
				"ca":  data.Category,
				"na":  data.Title,
				"au":  userId,
				"pr":  data.Preview,
				"te":  data.Content,
				"dcr": time.Now(),
				"dch": time.Now()}); err != nil {
			log.Printf("SaveRecord: %s\n", err.Error())
			return false
		} else {
			return true
		}
	} else {
		recordId := 0
		// check author and id from form
		// only author can edit record
		if err := q.Get(&recordId, `SELECT id FROM blog WHERE author = ? AND id = ?`, userId, data.Id); err != nil {
			log.Printf("SaveRecord: %s\n", err.Error())
		} else {
			if recordId > 0 {
				if _, err = q.NamedExec(`UPDATE blog SET name = :ti, category = :ca, preview = :pr, text = :te, datechanged = :dc
					WHERE id = :id`,
					map[string]interface{}{
						"ti": data.Title,
						"ca": data.Category,
						"pr": data.Preview,
						"te": data.Content,
						"dc": time.Now(),
						"id": recordId}); err != nil {
					log.Printf("SaveRecord: %s\n", err.Error())
					return false
				} else {
					return true
				}
			}
		}
	}
	return false
}

func DeleteRecord(userId int, data FormRecord) bool {
	q := models.GetConnection()
	recordId := 0
	// check author and id from form
	// only author can delete record
	if err := q.Get(&recordId, `SELECT id FROM blog WHERE author = ? AND id = ?`, userId, data.Id); err != nil {
		log.Printf("DeleteRecord: %s\n", err.Error())
	} else {
		if recordId > 0 {
			if _, err = q.NamedExec(`DELETE FROM blog WHERE id = :id`,
				map[string]interface{}{
					"id": recordId}); err != nil {
				log.Printf("DeleteRecord: %s\n", err.Error())
				return false
			} else {
				return true
			}
		}
	}
	return false
}

func DeleteCategory(data FormCategory) bool {
	q := models.GetConnection()
	if _, err := q.NamedExec(`DELETE FROM category WHERE id = :id`,
		map[string]interface{}{
			"id": data.Id}); err != nil {
		log.Printf("DeleteCategory: %s\n", err.Error())
		return false
	}
	return true
}
