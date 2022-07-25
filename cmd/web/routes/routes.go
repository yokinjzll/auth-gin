package routes

import (
	"github.com/gin-gonic/gin"

	controllers "gin-chat/cmd/web/controllers"
)

func PublicRoutes(g *gin.RouterGroup) {
	g.GET("/login", controllers.LoginGetHandler())
	g.POST("/login", controllers.LoginPostHandler())
	g.GET("/", controllers.IndexGetHandler())
	g.GET("/registration", controllers.RegistrationGetHandler())
	g.POST("/registration", controllers.RegistrationPostHandler())
}

func PrivateRoutes(g *gin.RouterGroup) {
	g.GET("/dashboard", controllers.DashboardGetHandler())
	g.GET("/logout", controllers.LogoutGetHandler())
	g.GET("/profile/:userid", controllers.ProfileGetHandler())
	g.GET("/profile/edit", controllers.ProfileEditGetHandler())
	g.POST("/profile/edit", controllers.ProfileEditPostHandler())
}
