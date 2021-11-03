package models

import u "bd_admin/utils"

type Teacher struct {
	AccountID     uint   `gorm:"primary_key"`
	FirstName  	  string `json:"first_name"`
	LastName 	  string `json:"last_name"`
}

func (teacher *Teacher) Create() map[string] interface{} {

	GetDB().Create(teacher)

	resp := u.Message(true, "success")
	resp["respond"] = teacher
	return resp
}