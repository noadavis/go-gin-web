package models

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var conn *sqlx.DB
var props ConfigManager

func InitDb() bool {
	var err error
	params := getConnectionString()
	// fmt.Println(fmt.Sprintf("db connection params: %s", params))
	conn, err = sqlx.Connect("mysql", params)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func CheckConnection() bool {
	if err := conn.Ping(); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func GetConnection() *sqlx.DB {
	return conn
}

func getConnectionString() string {
	var conf = props.GetProps().DbConf
	return fmt.Sprintf("%s:%s@(%s:%d)/%s", conf.User, conf.Password, conf.Address, conf.Port, conf.Name)
}
