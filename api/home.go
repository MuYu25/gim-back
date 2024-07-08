package api

import (
	"math/rand"
	"net/http"
	"project/utils/errmsg"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取活跃日期，待完成，未使用真正的数据
func GetHomeData(c *gin.Context) {
	queryParam := c.Query("date")
	dateInt, err := strconv.Atoi(queryParam)
	if queryParam == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
		})
		c.Abort()
		return
	}
	data := make(map[string]interface{})
	data["visits"] = 1000
	data["registrations"] = 10000
	data["downloads"] = 100000
	data["sales"] = 1000000
	data["recentVisits"] = make([][]int, 3)
	for i, _ := range data["recentVisits"].([][]int) {
		data["recentVisits"].([][]int)[i] = make([]int, dateInt)
		for j, _ := range data["recentVisits"].([][]int)[i] {
			data["recentVisits"].([][]int)[i][j] = rand.Intn(901) + 100
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": errmsg.SUCCESS,
		"data":    data,
	})
}
