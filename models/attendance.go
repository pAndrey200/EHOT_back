package models

import (
	u "bd_admin/utils"
	"fmt"
)

type Attendance struct {
	StudentId uint   `json:"student_id"`
	SubId     uint   `json:"sub_id"`
	Date      string `gorm:"type:date" json:"date"`
	Attend    bool   `json:"attendance"`
}

func (attend *Attendance) Create() map[string]interface{} {
	GetDB().Create(attend)
	resp := u.Message(true, "success")
	resp["respond"] = attend
	return resp
}

func GetAttendance(user uint) []*Attendance {

	attend := make([]*Attendance, 0)
	err := GetDB().Table("attendances").Where("student_id = ?", user).Find(&attend).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return attend
}
