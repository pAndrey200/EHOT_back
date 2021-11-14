package models

import (
	u "bd_admin/utils"
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