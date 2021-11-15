package models

import (
	u "bd_admin/utils"
	"fmt"
	"github.com/golang-sql/civil"
)

type Attendance struct {
	StudentId       uint  `gorm:"student_id"`
	SubId   		uint `json:"sub_id"`
	Date 			civil.Date `json:"date"`
	Attend   		bool `json:"attendance"`
}

func (attend *Attendance) Create() map[string] interface{} {
	GetDB().Create(attend)
	resp := u.Message(true, "success")
	resp["respond"] = attend
	return resp
}

func GetAttendance(user uint) []*Attendance {

	attend := make([]*Attendance, 0)
	err := GetDB().Table("attendances").Where("user_id = ?", user).Find(&attend).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return attend
}