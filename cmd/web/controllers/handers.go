package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	globals "gin-chat/cmd/web/globals"
	helpers "gin-chat/cmd/web/helpers"
	"gin-chat/pkg/models"
	"gin-chat/pkg/models/mysql"
)

func LoginGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		if user != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"content": "Please logout first",
				"user":    user,
			})
			return
		}
		c.HTML(http.StatusOK, "login.html", gin.H{
			"content": "",
			"user":    user,
		})
	}
}

func LoginPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		// user_id := session.Get(globals.UserID)
		// user_id := c.Param("user_id")
		if user != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"content": "Please logout first"})
			return
		}

		username := c.PostForm("username")
		password := c.PostForm("password")

		if helpers.EmptyUserPass(username, password) {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"content": "Parameters can't be empty"})
			return
		}

		if !helpers.CheckUserPass(username, password) {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"content": "Incorrect username or password"})
			return
		}

		db, err, _ := helpers.OpenDB(globals.Dsn)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"content": "Database's trables"})
			return
		}
		defer db.Close()

		usermodel := mysql.UserModel{DB: db}
		user_id, err := usermodel.Auth(username, password)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"content": "Ошибка при попытке авторизации."})
			return
		}

		// str_user_id := fmt.Sprintf("%v", user_id)

		// dbUserDetail, err, _ := helpers.OpenDB(globals.Dsn)
		// if err != nil {
		// 	log.Println(err)
		// 	c.HTML(http.StatusUnauthorized, "login.html", gin.H{"content": "Ошибка при открытии бд.", "user": user})
		// 	return
		// }
		// defer dbUserDetail.Close()

		// userdetailsmodel := mysql.UserDetailModel{DB: dbUserDetail}
		// user_details, err := userdetailsmodel.Get(str_user_id)
		// if err != nil {
		// 	log.Println(err)
		// 	c.HTML(http.StatusUnauthorized, "login.html", gin.H{"content": "Ошибка получения данный из бд.", "user": user})
		// 	return
		// }

		session.Set(globals.Userkey, username)
		session.Set(globals.UserID, user_id)
		// session.Set(globals.UserDetail, user_details)
		if err := session.Save(); err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{"content": "Failed to save session", "user": user})
			return
		}

		c.Redirect(http.StatusMovedPermanently, "/dashboard")
	}
}

func LogoutGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		// user_id := session.Get(globals.UserID)
		// user_details := session.Get(globals.UserDetail)
		log.Println("logging out user:", user)
		if user == nil {
			log.Println("Invalid session token")
			return
		}
		// session.Delete(globals.Userkey)
		// session.Delete(globals.UserID)
		// session.Delete(globals.UserDetail)
		session.Clear()
		session.Options(sessions.Options{MaxAge: -1})
		if err := session.Save(); err != nil {
			log.Println("Failed to save session:", err)
			return
		}

		c.Redirect(http.StatusMovedPermanently, "/")
	}
}

func IndexGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		user_id := session.Get(globals.UserID)
		user_details := session.Get(globals.UserDetail)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content":      "This is an index page...",
			"user":         user,
			"user_id":      user_id,
			"user_details": user_details,
		})
	}
}

func DashboardGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		user_id := session.Get(globals.UserID)
		str_user_id := fmt.Sprintf("%v", user_id)

		dbUserDetail, err, _ := helpers.OpenDB(globals.Dsn)
		if err != nil {
			log.Println(err)
		}
		defer dbUserDetail.Close()

		userdetailsmodel := mysql.UserDetailModel{DB: dbUserDetail}
		user_details, err := userdetailsmodel.Get(str_user_id)
		if err != nil {
			c.HTML(http.StatusBadRequest, "dashboard.html", gin.H{"content": "Не удалость загрузить данные."})
			return
		}

		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"content":      "This is a dashboard",
			"user":         user,
			"user_id":      str_user_id,
			"user_details": user_details,
		})
	}
}

func RegisterGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		if user != nil {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"content": "Please logout first",
				"user":    user,
			})
			return
		}
		c.HTML(http.StatusOK, "register.html", gin.H{
			"content": "",
			"user":    user,
		})
	}
}

func RegisterPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		// user_id := session.Get(globals.UserID)
		if user != nil {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"content": "Please logout first", "user": user})
			return
		}

		username := c.PostForm("username")
		password := c.PostForm("password")
		password_confirm := c.PostForm("password-confirm")
		first_name := c.PostForm("first-name")
		last_name := c.PostForm("last-name")
		dob := c.PostForm("date-of-birth")

		log.Println(username, password, password_confirm, first_name, last_name, dob)
		dbUser, err, _ := helpers.OpenDB(globals.Dsn)
		if err != nil {
			log.Println(err)
		}
		defer dbUser.Close()

		usermodel := mysql.UserModel{DB: dbUser}
		user_id, err := usermodel.Insert(username, password)
		if err != nil {
			log.Println(err)
		}

		// dobtime, err := time.Parse("2006-01-02", dob)
		nowtime := time.Now()
		userstructsend := models.UserDetails{UserID: user_id, First_name: first_name, Last_name: last_name, Dob: dob, Created_at: nowtime}

		dbUserDetail, err, _ := helpers.OpenDB(globals.Dsn)
		if err != nil {
			log.Println(err)
		}
		defer dbUserDetail.Close()

		userdetailsmodel := mysql.UserDetailModel{DB: dbUserDetail}
		usersendDetailID, err := userdetailsmodel.Insert(userstructsend)
		if err != nil {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"content": "Не удалость загрузить данные."})
			return
		}

		// user_details, err := userdetailsmodel.Get(usersendDetailID)
		// if err != nil {
		// 	log.Println(err)
		// }

		if helpers.EmptyUserDetails(username, password, password_confirm, first_name, last_name, dob) {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"content": "Parameters can't be empty"})
			return
		}

		if !helpers.EqualPasswords(password, password_confirm) {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"content": "Passwords are not equals!"})
			return
		}

		// if !helpers.CheckUserPass(username, password) {
		// 	c.HTML(http.StatusUnauthorized, "login.html", gin.H{"content": "Incorrect username or password"})
		// 	return
		// }

		// in future add more validation on check username and password

		session.Set(globals.Userkey, username)
		session.Set(globals.UserID, usersendDetailID)
		session.Set(globals.UserDetail, userstructsend)
		// session.Set(globals.UserDetail, user_details)
		if err := session.Save(); err != nil {
			c.HTML(http.StatusInternalServerError, "register.html", gin.H{"content": "Failed to save session"})
			return
		}

		// c.HTML(http.StatusInternalServerError, "register.html", gin.H{"content": "Failed to save session"})
		c.Redirect(http.StatusMovedPermanently, "/dashboard")
	}
}

func ProfileGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		// user_id := session.Get(globals.UserID)
		userid := c.Param("userid")
		log.Println(userid)
		// session.Set(globals.UserID, user_id)
		// useridsession := session.Get(globals.UserID)

		dbUserDetail, err, _ := helpers.OpenDB(globals.Dsn)
		if err != nil {
			log.Println(err)
		}
		userdetailsmodel := mysql.UserDetailModel{DB: dbUserDetail}
		user_details, err := userdetailsmodel.Get(userid)
		if err != nil {
			log.Println(err)
		}
		// user_det := session.Get(globals.UserDetail)
		c.HTML(http.StatusOK, "profile.html", gin.H{
			"content":      "This is a your profile.",
			"user":         user,
			"user_id":      userid,
			"user_details": user_details,
		})
	}
}

func ProfileEditGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		user_id := session.Get(globals.UserID)
		str_user_id := fmt.Sprintf("%v", user_id)

		dbUserDetail, err, _ := helpers.OpenDB(globals.Dsn)
		if err != nil {
			log.Println("[error]", err, "db not open")
		}
		defer dbUserDetail.Close()

		userdetailsmodel := mysql.UserDetailModel{DB: dbUserDetail}
		user_details, err := userdetailsmodel.Get(str_user_id)
		if err != nil {
			c.HTML(http.StatusBadRequest, "profile_edit.html", gin.H{
				"content": "Не удалость загрузить данные."})
			return
		}

		// dob := c.PostForm("date-of-birth")
		// log.Println(dob)
		// user_details.Dob = time.Parse("", user_details.Dob)

		c.HTML(http.StatusOK, "profile_edit.html", gin.H{
			"content":      "This is a your profile editor.",
			"user":         user,
			"user_id":      user_id,
			"user_details": user_details,
		})
	}
}

func ProfileEditPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		user_id := session.Get(globals.UserID)
		str_user_id := fmt.Sprintf("%v", user_id)

		int_user_id, err := strconv.Atoi(str_user_id)
		if err != nil {
			log.Println(err)
		}
		// user_det := session.Get(globals.UserDetail)

		first_name := c.PostForm("first-name")
		last_name := c.PostForm("last-name")
		dob := c.PostForm("date-of-birth")

		log.Println("edit dob to ", dob)

		dbUserDetail, err, _ := helpers.OpenDB(globals.Dsn)
		if err != nil {
			log.Println(err)
		}
		defer dbUserDetail.Close()
		userdetailsmodel := mysql.UserDetailModel{DB: dbUserDetail}

		userstructsend := models.UserDetails{UserID: int_user_id, First_name: first_name, Last_name: last_name, Dob: dob}
		user_details_old, err := userdetailsmodel.Get(str_user_id)
		if err != nil {
			c.HTML(http.StatusBadRequest, "profile_edit.html", gin.H{
				"content":      "Не удалость получить данные.",
				"user":         user,
				"user_id":      user_id,
				"user_details": user_details_old})
			return
		}

		usersendDetailID := userdetailsmodel.SetByID(str_user_id, userstructsend)
		if !usersendDetailID {
			c.HTML(http.StatusBadRequest, "profile_edit.html", gin.H{
				"content":      "Не удалость загрузить данные.",
				"user":         user,
				"user_id":      int_user_id,
				"user_details": user_details_old})
			return
		}

		user_details, err := userdetailsmodel.Get(str_user_id)
		if err != nil {
			c.HTML(http.StatusBadRequest, "profile_edit.html", gin.H{
				"content":      "Не удалость получить данные.",
				"user":         user,
				"user_id":      user_id,
				"user_details": user_details_old})
			return
		}

		log.Println(user_details.Dob)
		// session.Set(globals.UserDetail, user_details)

		// c.HTML(http.StatusOK, "profile_edit.html", gin.H{
		// 	"content":      "This is a your profile editor.",
		// 	"user":         user,
		// 	"user_id":      user_id,
		// 	"user_details": user_details,
		// })

		// var url_redirect string = "profile/" + int_user_id

		c.Redirect(http.StatusMovedPermanently, "/profile/"+str_user_id)
	}
}
