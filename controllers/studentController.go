package controllers

import (
	"bd_admin/models"
	u "bd_admin/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

var CreateStudent = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint) //Получение идентификатора пользователя, отправившего запрос
	student := &models.Student{}
	//fmt.Println(user)
	err := json.NewDecoder(r.Body).Decode(student)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body "+err.Error()))
		return
	}

	student.AccountID = user
	resp := student.Create()
	u.Respond(w, resp)
}

var GetStudentAttendance = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	data := models.GetAttendance(user)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var GetInfo = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	acc := models.Account{}
	err := models.GetDB().Table("accounts").Where("ID = ?", user).First(&acc).Error
	if err != nil {
		u.Respond(w, u.Message(false, "Error while exec query "+err.Error()))
		return
	}
	resp := u.Message(true, "success")
	if acc.Role == "student" {
		st := models.Student{}
		err := models.GetDB().Table("students").Where("account_id = ?", user).First(&st).Error
		if err != nil {
			u.Respond(w, u.Message(false, "Error while exec query "+err.Error()))
			return
		}
		resp["data"] = st
	} else {
		t := models.Teacher{}
		err := models.GetDB().Table("teachers").Where("account_id = ?", user).First(&t).Error
		if err != nil {
			u.Respond(w, u.Message(false, "Error while exec query "+err.Error()))
			return
		}
		resp["data"] = t
	}
	u.Respond(w, resp)
}

var SentQRCode = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	str := strconv.Itoa(int(user)) + " true"
	resp := u.Message(true, "success")
	resp["token"] = str
	u.Respond(w, resp)
}
