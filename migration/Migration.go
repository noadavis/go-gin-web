package migration

import (
	"database/sql"
	"fmt"
	"go-gin-web/models"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var props models.ConfigManager

func InitMigration() {
	if len(os.Args) > 2 {
		if os.Args[1] != "migrate" {
			return
		}
		switch os.Args[2] {
		case "up":
			RunMigration(true)
			return
		case "down":
			RunMigration(false)
			return
		}
	}
	log.Println("Unknown arguments")
}

func RunMigration(up bool) {
	// connect to db
	conn, err := sql.Open("mysql", getConnectionString())
	if err != nil {
		log.Fatalf("could not connect to the MySQL database... %v", err)
	}
	// check db connection
	if err := conn.Ping(); err != nil {
		log.Fatalf("could not ping DB... %v", err)
	}
	// read migrations
	driver, err := mysql.WithInstance(conn, &mysql.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://./migration/migrations", "mysql", driver)
	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}
	// migrate
	if up {
		// apply migrations
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("An error occurred while syncing the database.. %v", err)
		}
		log.Println("Database migrated")
	} else {
		// get migration version
		if version, _, err := m.Version(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("An error occurred while read database version.. %v", err)
		} else {
			log.Printf("Migrate to version: %d\n", version-1)
			if version > 0 {
				// drop last migration
				if version > 1 {
					if err := m.Migrate(version - 1); err != nil && err != migrate.ErrNoChange {
						log.Fatalf("An error occurred while syncing the database.. %v", err)
					}
				} else {
					if err := m.Down(); err != nil && err != migrate.ErrNoChange {
						log.Fatalf("An error occurred while syncing the database.. %v", err)
					}
				}
				log.Println("Drop last database migration")
			} else {
				log.Println("Database migrations not exist")
			}
		}
	}
}

func getConnectionString() string {
	var conf = props.GetProps().DbConf
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?multiStatements=true", conf.User, conf.Password, conf.Address, conf.Port, conf.Name)
}
