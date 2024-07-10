package model

import (
	"project/utils/errmsg"
	"time"

	"gorm.io/gorm"
)

type UserData struct {
	gorm.Model
	CreatedAt time.Time `gorm:"index" json:"CreatedAt" validate:"required" label:"创建时间"`
	Type      string    `gorm:"type:varchar(20);not null" json:"type" validate:"required" label:"类型"`
	Phone     string    `gorm:"type:varchar(20);not null" json:"phone" validate:"required,min=10,max=20" label:"手机号"`
	Cuy       string    `gorm:"type:varchar(100);not null" json:"cuy" validate:"required" label:"CUY字段"`
	Cc        string    `gorm:"type:varchar(100);not null" json:"cc" validate:"required" label:"CC字段"`
	Data      string    `gorm:"type:Longtext" json:"data" validate:"required" label:"数据内容"`
	Channel   string    `gorm:"type:varchar(50);not null" json:"channel" validate:"required" label:"渠道"`
	State     string    `gorm:"type:varchar(20);not null" json:"state" validate:"required" label:"状态"`
	Belong    string    `gorm:"type:varchar(100);not null" json:"belong" label:"归属信息"`
	ToUserId  int       `gorm:"not null" json:"to_user_id" validate:"required" label:"所属用户ID"`
}

// CreateUserData 添加一个条新的记录
func CreateUserData(data *UserData) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// GetUserData 根据查询条件
func GetUserData(_conditions UserData) (userData []UserData, total int64) {
	query := db.Model(&userData)
	if _conditions.Type != "" {
		query = query.Where("type = ?", _conditions.Type)
	}
	if _conditions.Cuy != "" {
		query = query.Where("cuy = ?", _conditions.Cuy)
	}
	if _conditions.Cc != "" {
		query = query.Where("cc = ?", _conditions.Cc)
	}
	if _conditions.State != "" {
		query = query.Where("state = ?", _conditions.State)
	}
	if _conditions.Channel != "" {
		query = query.Where("channel = ?", _conditions.Channel)
	}
	if !_conditions.CreatedAt.IsZero() {
		query = query.Where("DATE(created_at) = ?", _conditions.CreatedAt.Format("2006-01-02"))
	}
	if err := query.Where("to_user_id = ?", _conditions.ToUserId).Find(&userData).Count(&total).Error; err != nil {
		return userData, 0
	}
	return
}

// GetDataFormIds 根据所拿到的id返回数据
func GetDataByIds(ids []int) (result []UserData, total int) {
	if err := db.Where("id IN ?", ids).Find(&result).Error; err != nil {
		return result, 0
	}
	return result, len(result)
}
