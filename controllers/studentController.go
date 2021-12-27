package controllers

import (
	"bd_admin/models"
	u "bd_admin/utils"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
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
	curTime := time.Now()
	fmt.Println(curTime)
	s := models.Schedule{}
	ss := models.Student{}
	err := models.GetDB().Table("students").Where("account_id = ?", user).First(&ss).Error
	if err != nil {
		u.Respond(w, u.Message(false, "Error while exec query "+err.Error()))
		return
	}
	err = models.GetDB().Table("schedules").Where("schedules.group = ? AND time > ? AND time < ? AND day = ?", ss.Group, curTime.Add(-40*time.Minute).Format("15:04:05"), curTime.Add(10*time.Minute).Format("15:04:05"), strings.ToLower(time.Now().Weekday().String())).First(&s).Error
	//err := models.GetDB().Raw("SELECT * FROM schedules WHERE schedules.group = 'fn11-33b' AND  (time >= (CURRENT_TIME - interval '40 minutes') OR time <= (CURRENT_TIME + interval '10 minutes'))").Find(&s).Error
	var str string
	if err == gorm.ErrRecordNotFound {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "there is no lesson now"))
		return
	} else {
		str = strconv.Itoa(int(user)) + " " + strconv.Itoa(int(s.SubID)) + " " + time.Now().Format("2006-01-02") + " " + "true"
	}
	resp := u.Message(true, "success")
	resp["token"] = str
	u.Respond(w, resp)
}
