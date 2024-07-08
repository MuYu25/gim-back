package api

import (
	"fmt"
	"net/http"
	"project/model"
	"project/utils/errmsg"
	"project/utils/validator"
	"strings"

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
	_conditions := make(map[string][]interface{})
	queryParams := c.Request.URL.Query()
	fmt.Println(queryParams)
	for param, conditions := range queryParams {
		for _, condition := range conditions {
			if len(strings.TrimSpace(condition)) > 0 {
				_conditions[param] = append(_conditions[param], condition)
			}
		}
	}
	data, total := model.GetUserData(_conditions)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  errmsg.SUCCESS,
			"message": errmsg.GetErrMsg(errmsg.SUCCESS),
			"data":    data,
			"total":   total,
		},
	)
}
