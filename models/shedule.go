package models

import (
	u "bd_admin/utils"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Schedule struct {
	SubID          uint `gorm:"primaryKey;autoIncrement:false"`
	TeacherID      uint
	TeacherName    string `json:"teacher_name"`
	TeacherSurname string `json:"teacher_surname"`
	Group          string `json:"group" gorm:"primaryKey;autoIncrement:false"`
	Time           string `gorm:"type:time" json:"time"`
	Day            string `json:"day"`
	SubName        string `json:"subject"`
}

func (schedule *Schedule) Create() map[string]interface{} {
	temp := &Teacher{}
	fmt.Println(schedule)
	err := GetDB().Table("teachers").Where("first_name = ? AND last_name = ?", schedule.TeacherName, schedule.TeacherSurname).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry")
	}
	if temp.AccountID == 0 {
		return u.Message(false, "There is no such teacher")
	}
	schedule.TeacherID = temp.AccountID
	GetDB().Create(schedule)

	resp := u.Message(true, "success")
	resp["respond"] = schedule
	return resp
}
