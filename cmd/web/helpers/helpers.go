package helpers

import (
	"log"
	"strings"

	"gin-chat/cmd/web/globals"
	"gin-chat/pkg/models/mysql"
)

func CheckUserPass(username, password string) bool {
	log.Println("checkUserPass user:", username)

	db, err, _ := OpenDB(globals.Dsn)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	usermodel := mysql.UserModel{DB: db}
	userid, err := usermodel.Auth(username, password)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("User is [", userid, "] log in.")
	return true
}

func EmptyUserPass(username, password string) bool {
	return strings.Trim(username, " ") == "" || strings.Trim(password, " ") == ""
}
