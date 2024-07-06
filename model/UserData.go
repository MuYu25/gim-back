package model

import (
	"fmt"
	"project/utils/errmsg"

	"gorm.io/gorm"
)

type UserData struct {
	gorm.Model
	Type     string `gorm:"type:varchar(20);not null" json:"type" validate:"required,type=string" label:"类型"`
	Phone    string `gorm:"type:varchar(20);not null" json:"phone" validate:"optional,min=10,max=20" label:"手机号"`
	Cuy      string `gorm:"type:varchar(100);not null" json:"cuy" validate:"required" label:"CUY字段"`
	Cc       string `gorm:"type:varchar(100);not null" json:"cc" validate:"required" label:"CC字段"`
	Data     string `gorm:"type:Longtext"; json:"data" validate:"required" label:"数据内容"`
	Channel  string `gorm:"type:varchar(50);not null" json:"channel" validate:"required" label:"渠道"`
	State    string `gorm:"type:varchar(20);not null" json:"state" validate:"required,min=1,max=3" label:"状态"`
	Belong   string `gorm:"type:varchar(100);not null" json:"belong" validate:"optional,type=string" label:"归属信息"`
	ToUserId uint   `gorm:"not null" json:"to_user_id" validate:"required,type=int" label:"所属用户ID"`
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
func GetUserData(_conditions map[string][]interface{}) (userData []UserData, total int64) {
	// query := db.Model(&userData)
	_db := db
	fmt.Println(_conditions)
	for param, conditions := range _conditions {
		for _, condition := range conditions {
			_db = _db.Where(param+" = ?", condition)
		}
	}
	_db.Find(&userData)
	_db.Model(&userData).Count(&total)
	return
}
