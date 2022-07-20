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

		usermodel := helpers.GetUserModel()
		user_id, err := usermodel.Auth(username, password)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"content": "Ошибка при попытке авторизации."})
			return
		}

		session.Set(globals.Userkey, username)
		session.Set(globals.UserID, user_id)
		if err := session.Save(); err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{"content": "Failed to save session"})
			return
		}

		c.Redirect(http.StatusMovedPermanently, "/dashboard")
	}
}

func LogoutGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		log.Println("logging out user:", user)
		if user == nil {
			log.Println("User not logged in")
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{"content": "User not logged in"})
			return
		}
		session.Clear()
		session.Options(sessions.Options{MaxAge: -1})
		if err := session.Save(); err != nil {
			log.Println("Failed to save session:", err)
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{"content": "Failed to save session"})
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
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "This is an index page...",
			"user":    user,
			"user_id": user_id,
		})
	}
}

func DashboardGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		user_id := session.Get(globals.UserID)
		str_user_id := fmt.Sprintf("%v", user_id)
		log.Println("user ud:", user_id, "str user id", str_user_id)

		userdetailsmodel := helpers.GetUserDetailsModel()
		user_details, err := userdetailsmodel.Get(str_user_id)
		if err != nil {
			log.Println(err)
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

		if user != nil {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"content": "Please logout first"})
			return
		}

		username := c.PostForm("username")
		password := c.PostForm("password")
		password_confirm := c.PostForm("password-confirm")
		first_name := c.PostForm("first-name")
		last_name := c.PostForm("last-name")
		dob := c.PostForm("date-of-birth")

		log.Println(username, password, password_confirm, first_name, last_name, dob)

		usermodel := helpers.GetUserModel()
		user_id, err := usermodel.Insert(username, password)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "register.html", gin.H{"content": "Не удалось создать пользователя."})
			return
		}

		nowtime := time.Now()
		userstructsend := models.UserDetails{UserID: user_id, First_name: first_name, Last_name: last_name, Dob: dob, Created_at: nowtime}

		userdetailsmodel := helpers.GetUserDetailsModel()
		usersendDetailID, err := userdetailsmodel.Insert(userstructsend)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"content": "Не удалость загрузить данные."})
			return
		}
		log.Println(usersendDetailID)

		if helpers.EmptyUserDetails(username, password, password_confirm, first_name, last_name, dob) {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"content": "Parameters can't be empty"})
			return
		}

		if !helpers.EqualPasswords(password, password_confirm) {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"content": "Passwords are not equals!"})
			return
		}

		c.Redirect(http.StatusMovedPermanently, "/")
	}
}

func ProfileGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		userid := c.Param("userid")
		log.Println(userid)

		userdetailsmodel := helpers.GetUserDetailsModel()
		user_details, err := userdetailsmodel.Get(userid)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "profile.html", gin.H{"content": "Не удалось загрузить данные пользователя."})
			return
		}

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

		userdetailsmodel := helpers.GetUserDetailsModel()
		user_details, err := userdetailsmodel.Get(str_user_id)
		if err != nil {
			c.HTML(http.StatusBadRequest, "profile_edit.html", gin.H{
				"content": "Не удалость загрузить данные."})
			return
		}

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

		first_name := c.PostForm("first-name")
		last_name := c.PostForm("last-name")
		dob := c.PostForm("date-of-birth")

		userdetailsmodel := helpers.GetUserDetailsModel()

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

		c.Redirect(http.StatusMovedPermanently, "/profile/"+str_user_id)
	}
}
