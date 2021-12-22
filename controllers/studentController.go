package controllers

import (
	"bd_admin/models"
	u "bd_admin/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

var GetQRCode = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint)
	curTime := time.Now()
	s := models.Schedule{}
	fmt.Println(curTime.Add(-40 * time.Minute).Format("15:04:05"))
	//err := models.GetDB().Table("schedules").Where("schedules.group = 'fn11-33b' AND (time > ? OR time < ?", curTime.Add(-40*time.Minute).Format("15:04:05"), curTime.Add(10*time.Minute).Format("15:04:05")).Find(&s).Error
	err := models.GetDB().Raw("SELECT * FROM schedules WHERE schedules.group = 'fn11-33b' AND (time > curtime OR time < ?", curTime.Add(-40*time.Minute).Format("15:04:05"), curTime.Add(10*time.Minute)).First(&s).Error

	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Error while exec query"))
		return
	}
	//data := models.GetAttendance(user)
	resp := u.Message(true, "success")
	resp["sub_id"] = s.SubID
	resp["user"] = user
	//resp["data"] = data
	u.Respond(w, resp)
}
