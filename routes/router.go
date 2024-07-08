package routes

import (
	"project/api"
	"project/middleware"
	"project/utils"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	r := gin.New()
	// 设置信任网络 []string
	// nil 为不计算，避免性能消耗，上线应当设置
	_ = r.SetTrustedProxies(nil)

	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	auth := r.Group("api")
	auth.Use(middleware.JwtToken())
	{
		// 用户模块的路由接口
		// auth.GET("admin/users", api.GetUsers)
		// auth.PUT("user/:id", api.EditUser)
		// auth.DELETE("user/:id", api.DeleteUser)
		auth.GET("home", api.GetHomeData)
		auth.GET("data", api.GetUserData)
		auth.POST("download", api.DownloadData)
	}

	router := r.Group("api")
	{
		// 用户信息模块
		router.POST("user/add", api.AddUser)
		router.GET("user/:id", api.GetUserInfo)
		router.GET("users", api.GetUsers)
		// 登录控制模块
		router.POST("login", api.Login)
		// router.POST("loginfront", api.LoginFront)

		router.POST("upload", api.AddData)
	}
	_ = r.Run(utils.HttpPort)
}
