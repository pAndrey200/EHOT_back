package models

import (
	u "bd_admin/utils"
	_ "github.com/dgrijalva/jwt-go"
	_ "os"
	_ "strings"
)

type Student struct {
	AccountID uint   `gorm:"primary_key"`
	FirstName string `json:"first_name"`
	Group     string `json:"group"`
	Year      uint   `json:"year"`
}

func (student *Student) Create() map[string]interface{} {

	//if resp, ok := contact.Validate(); !ok {
	//	return resp
	//}

	GetDB().Create(student)

	resp := u.Message(true, "success")
	resp["respond"] = student
	return resp
}
