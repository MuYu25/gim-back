package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"project/model"
	"project/utils"
	"project/utils/errmsg"
	"project/utils/validator"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// 根据条件下载文件
func DownloadData(c *gin.Context) {
	type data struct {
		Ids []int `json:"fileId" validate:"required" label:"下载条件"`
	}
	var ids data
	_ = c.ShouldBindJSON(&ids)
	// 验证参数
	msg, validCode := validator.Validate(&ids)
	if validCode != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  validCode,
			"message": msg,
		})
		c.Abort()
		return
	}

	// 获取数据
	results, total := model.GetDataFormIds(ids.Ids)
	if total == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
		})
		c.Abort()
		return
	}

	// 获取token
	tokenString := c.GetHeader("Authorization")
	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.JwtKey), nil
	}); err != nil {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"status":  errmsg.ERROR,
				"message": errmsg.GetErrMsg(errmsg.ERROR),
			},
		)
	}
	username := claims["username"].(string)
	bytes := make([]byte, 3)
	rand.Read(bytes)
	filename := fmt.Sprintf("%s_%s_%s.txt", username, time.Now().Format("2006-01-02"), base64.StdEncoding.EncodeToString(bytes))
	filepath := utils.FilePath + filename
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// 文件不存在，创建文件并写入数据

		file, _ := os.Create(filepath)
		defer file.Close()
		for _, result := range results {
			data, _ := json.Marshal(result)
			file.Write(data)
			file.Write([]byte("\n"))
		}
	}
	// c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filepath)
}
