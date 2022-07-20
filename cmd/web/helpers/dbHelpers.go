package helpers

import (
	"database/sql"
	"gin-chat/pkg/models/mysql"
	"log"
	"strings"
)

var dbUserModel *mysql.UserModel
var dbUserDetailsModel *mysql.UserDetailsModel

func EmptyUserDetails(username, password, password_confirm, first_name, last_name, dob string) bool {

	return strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" || strings.Trim(password_confirm, " ") == "" || strings.Trim(first_name, " ") == "" || strings.Trim(last_name, " ") == "" || strings.Trim(dob, " ") == ""
}

func EqualPasswords(pass1, pass2 string) bool {
	if pass1 == pass2 {
		return true
	}
	return false
}

func OpenDB(dsn string) (*sql.DB, error, bool) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err, false
	}
	if err = db.Ping(); err != nil {
		return nil, err, false
	}
	return db, nil, true
}

func SetUserModel(db *sql.DB, err error) {
	dbUserModel = &mysql.UserModel{DB: db}
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func SetUserDetailsModel(db *sql.DB, err error) {
	dbUserDetailsModel = &mysql.UserDetailsModel{DB: db}
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func GetUserModel() *mysql.UserModel {
	return dbUserModel
}

func GetUserDetailsModel() *mysql.UserDetailsModel {
	return dbUserDetailsModel
}
