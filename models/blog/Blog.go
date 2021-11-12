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
}

func GetCategories(auth bool) []BlogCategory {
	q := models.GetConnection()
	categories := []BlogCategory{}
	showAll := ""
	if !auth {
		showAll = "WHERE auth = 0"
	}
	err := q.Select(&categories, fmt.Sprintf(`SELECT id, name, alias
		FROM category %s ORDER BY id`, showAll))
	if err != nil {
		log.Printf("GetCategories: %s", err.Error())
	}
	return categories
}

func GetCategory(category int, auth bool) []BlogRecord {
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
	recordId := 0
	// check author and id from form
	// only author can edit record
	if err := q.Get(&recordId, `SELECT id FROM blog WHERE author = ? AND id = ?`, userId, data.Id); err != nil {
		log.Printf("SaveRecord: %s\n", err.Error())
	} else {
		if recordId > 0 {
			if _, err = q.NamedExec(`UPDATE blog SET name = :ti, category = :ca, preview = :re, text = :te, datechanged = :dc
				WHERE id = :id`,
				map[string]interface{}{
					"ti": data.Title,
					"ca": data.Category,
					"re": data.Preview,
					"te": data.Content,
					"dc": time.Now(),
					"id": recordId}); err != nil {
				log.Printf("EditUser [user]: %s\n", err.Error())
				return false
			} else {
				return true
			}
		}
	}
	return false
}
