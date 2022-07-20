package main

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	globals "gin-chat/cmd/web/globals"
	helpers "gin-chat/cmd/web/helpers"
	middleware "gin-chat/cmd/web/middleware"
	routes "gin-chat/cmd/web/routes"

	// mysql "gin-chat/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

// type DbModels struct {
// 	UserModel        *sql.DB
// 	UserDetailsModel *sql.DB
// }

func main() {

	db, err, _ := helpers.OpenDB(globals.Dsn)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	helpers.SetUserModel(db, err)
	helpers.SetUserDetailsModel(db, err)

	router := gin.Default()

	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.html")

	router.Use(sessions.Sessions("session", cookie.NewStore(globals.Secret)))

	public := router.Group("/")
	routes.PublicRoutes(public)

	private := router.Group("/")
	private.Use(middleware.AuthRequired)
	routes.PrivateRoutes(private)

	router.Run("localhost:8080")
}
