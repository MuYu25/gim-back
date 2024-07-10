package api

import (
	"net/http"
	"project/middleware"
	"project/model"
	"project/utils/errmsg"
	"project/utils/validator"
	"time"

	"github.com/gin-gonic/gin"
)

// 添加数据
func AddData(c *gin.Context) {
	var data model.UserData
	_ = c.ShouldBindJSON(&data)
	// var code int
	// fmt.Println(data)
	msg, validcode := validator.Validate(&data)
	if validcode != errmsg.SUCCESS {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  validcode,
				"message": msg,
			},
		)
		c.Abort()
		return
	}
	model.CreateUserData(&data)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  errmsg.SUCCESS,
			"message": errmsg.GetErrMsg(errmsg.SUCCESS),
		},
	)
}

// 查询用户的数据 /api/data
func GetUserData(c *gin.Context) {
	var conditions model.UserData
	// conditions := make(model.UserData)
	conditions = model.UserData{
		Cc:       c.Query("cc"),
		Type:     c.Query("type"),
		Channel:  c.Query("channel"),
		State:    c.Query("state"),
		ToUserId: int(middleware.NewJWT().GetTokenValues(c.GetHeader("Authorization"), "uid").(float64)),
	}
	if c.Query("CreatedAt") != "" {
		conditions.CreatedAt, _ = time.Parse("2006-1-2", c.Query("CreatedAt"))
	}
	results, total := model.GetUserData(conditions)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  errmsg.SUCCESS,
			"message": errmsg.GetErrMsg(errmsg.SUCCESS),
			"data":    results,
			"total":   total,
		},
	)
}
