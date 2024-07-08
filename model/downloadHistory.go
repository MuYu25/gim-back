package model

import (
	"project/utils/errmsg"

	"gorm.io/gorm"
)

type DownloadHistory struct {
	gorm.Model
	UserId   int    `gorm:"type:int;not null" json:"user_id" validate:"required,type=string" label:"所属用户id"`
	FileName string `gorm:"type:varchar(100); not null" josn:"file_name" validate:"optional" label:"文件名"`
}

func CreateDownloadHistory(data *DownloadHistory) int {
	if err := db.Create(&data).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

func GetDownloadHistory(id int) ([]DownloadHistory, int) {
	var downloadHistory []DownloadHistory
	if err := db.Where("user_id = ?", id).Find(&downloadHistory).Error; err != nil {
		return downloadHistory, 0
	}
	return downloadHistory, len(downloadHistory)
}
