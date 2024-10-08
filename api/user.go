package api

import (
	"net/http"
	"project/model"
	"project/utils"
	"project/utils/errmsg"
	"strconv"

	"project/utils/validator"

	"github.com/gin-gonic/gin"
)

// AddUser 添加用户
func AddUser(c *gin.Context) {
	type UserRequest struct {
		model.User
		InviteCode string `json:"invite_code" validate:"required" label:"邀请码"`
	}
	var data UserRequest
	var msg string
	var validCode int
	_ = c.ShouldBindJSON(&data)
	msg, validCode = validator.Validate(&data)
	if validCode != errmsg.SUCCESS {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"status":  validCode,
				"message": msg,
			},
		)
		c.Abort()
		return
	}

	// 写死的邀请码
	if data.InviteCode == "" {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  500,
				"message": "邀请码不能为空",
			},
		)
		c.Abort()
		return
	} else if data.InviteCode != utils.InvitationCode {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"status":  500,
				"message": "邀请码错误",
			},
		)
		c.Abort()
		return
	}

	code := model.CheckUser(data.Username)
	if code == errmsg.SUCCESS {
		model.CreateUser(&data.User)
	}

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// GetUserInfo 查询单个用户
func GetUserInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var maps = make(map[string]interface{})
	data, code := model.GetUser(id)
	maps["username"] = data.Username
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    maps,
			"total":   1,
			"message": errmsg.GetErrMsg(code),
		},
	)

}

// GetUsers 查询用户列表
func GetUsers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	username := c.Query("username")

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	data, total := model.GetUsers(username, pageSize, pageNum)

	code := errmsg.SUCCESS
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   total,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// EditUser 编辑用户
func EditUser(c *gin.Context) {
	var data model.User
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)

	code := model.CheckUpUser(id, data.Username)
	if code == errmsg.SUCCESS {
		model.EditUser(id, &data)
	}

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// ChangeUserPassword 修改密码
func ChangeUserPassword(c *gin.Context) {
	var data model.User
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)

	code := model.ChangePassword(id, &data)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	code := model.DeleteUser(id)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}
