package model

import (
	"fmt"
	"os"
	"project/utils"
	"project/utils/errmsg"
	"time"

	"gorm.io/gorm"
)

type DownloadHistory struct {
	gorm.Model
	UserId   int    `gorm:"type:int;not null" json:"user_id" validate:"required" label:"所属用户id"`
	FileName string `gorm:"type:varchar(100); not null" json:"file_name" label:"文件名"`
}

// 创建新的下载历史
func CreateDownloadHistory(data *DownloadHistory) int {
	if err := db.Create(&data).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 根据userid返回相对应的下载历史
func GetDownloadHistoryByUserIdAndDate(created_at time.Time, id int) ([]DownloadHistory, int) {
	var downloadHistory []DownloadHistory
	if !created_at.IsZero() {
		db.Where("DATE(created_at) = ? ", created_at.Format("2006-1-2"))
	}
	if err := db.Where("user_id = ?", id).Find(&downloadHistory).Error; err != nil {
		return downloadHistory, 0
	}
	return downloadHistory, len(downloadHistory)
}

// 根据记录id查询下载历史
func GetDownloadHistoryById(id int) (DownloadHistory, int) {
	var downloadHistory DownloadHistory
	if err := db.Where("id = ?", id).Find(&downloadHistory).Error; err != nil {
		return DownloadHistory{}, errmsg.ERROR
	}
	return downloadHistory, errmsg.SUCCESS
}

// 根据id修改下载历史
func UpdateFileById(id int, newFilename string) int {
	var downloadHistory DownloadHistory
	if err := db.Where("id = ?", id).Find(&downloadHistory).Error; err != nil {
		fmt.Println("查询记录不存在")
		return errmsg.ERROR
	}
	newFIlePath := utils.FilePath + newFilename
	if _, err := os.Stat(newFIlePath); !os.IsNotExist(err) {
		return errmsg.ERROR_FILE_EXIST
	}
	go func() {
		os.Rename(utils.FilePath+downloadHistory.FileName, newFIlePath)
		downloadHistory.FileName = newFilename
		db.Save(&downloadHistory)
	}()
	return errmsg.SUCCESS
}
