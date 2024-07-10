package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"project/middleware"
	"project/model"
	"project/utils"
	"project/utils/errmsg"
	"project/utils/validator"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 根据条件下载文件
func DownloadData(c *gin.Context) {
	type data struct {
		Ids []int `json:"fileId" validate:"required,min=1" label:"下载条件"`
	}
	var ids data
	_ = c.ShouldBindJSON(&ids)
	// 验证参数
	msg, validCode := validator.Validate(&ids)
	if validCode != errmsg.SUCCESS {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  validCode,
			"message": msg,
		})
		c.Abort()
		return
	}

	// 获取数据
	results, total := model.GetDataByIds(ids.Ids)
	if total == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
		})
		c.Abort()
		return
	}

	// 获取token 并且解析拿到token中的username值
	tokenString := c.GetHeader("Authorization")
	j := middleware.NewJWT()
	username := j.GetTokenValues(tokenString, "username").(string)
	// 去创建下载历史记录
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
	go func() {
		uid := int(j.GetTokenValues(tokenString, "uid").(float64))
		if model.CreateDownloadHistory(&model.DownloadHistory{UserId: uid, FileName: filename}) == errmsg.ERROR {
			fmt.Println("创建下载历史记录失败, 可能是文件名冲突")
		}
	}()
	// c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filepath)
}

// 获取当前用户的所有下载历史
func GetDownloadHistory(c *gin.Context) {
	var createTime time.Time
	var err error
	if c.Query("CreatedAt") != "" {
		if createTime, err = time.Parse("2006-1-2", c.Query("CreatedAt")); err != nil {
			c.JSON(
				http.StatusInternalServerError, gin.H{
					"status":  errmsg.ERROR,
					"message": "时间格式错误",
				},
			)
			c.Abort()
			return
		}
	}

	j := middleware.NewJWT()
	uid := int(j.GetTokenValues(c.GetHeader("Authorization"), "uid").(float64))
	results, total := model.GetDownloadHistoryByUserIdAndDate(createTime, uid)
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"data":    results,
		"total":   total,
		"message": errmsg.GetErrMsg(errmsg.SUCCESS),
	})
}

// 根据文件下载历史的Id修改文件名
func UpdateDownloadHistory(c *gin.Context) {
	type condition struct {
		Id       int    `json:"id" validate:"required,min=1" label:"下载历史Id"`
		FileName string `json:"fileName" validate:"required,min=1" label:"新的文件名"`
	}
	var data condition
	_ = c.ShouldBindJSON(&data)
	fmt.Println(data)
	msg, validCode := validator.Validate(&data)
	if validCode != errmsg.SUCCESS {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  validCode,
			"message": msg,
		})
		c.Abort()
		return
	}
	if model.UpdateFileById(data.Id, data.FileName) != errmsg.SUCCESS {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR_FILE_EXIST),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": errmsg.GetErrMsg(errmsg.SUCCESS),
	})
}

// 再次下载
func DownloadById(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "请传入整数类型",
		})
		c.Abort()
		return
	}
	downloadHistory, code := model.GetDownloadHistoryById(id)
	if code != errmsg.SUCCESS {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"status":  errmsg.ERROR,
				"message": errmsg.GetErrMsg(errmsg.ERROR),
			},
		)
	}
	c.Header("Content-Disposition", "attachment; filename="+downloadHistory.FileName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(utils.FilePath + downloadHistory.FileName)
}
