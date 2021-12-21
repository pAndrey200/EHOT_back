package controllers

import (
	"bd_admin/models"
	u "bd_admin/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var CreateTeacher = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint) //Получение идентификатора пользователя, отправившего запрос
	teacher := &models.Teacher{}

	err := json.NewDecoder(r.Body).Decode(teacher)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	teacher.AccountID = user
	resp := teacher.Create()
	u.Respond(w, resp)
}

var CreateSub = func(w http.ResponseWriter, r *http.Request) {
	schedule := &models.Schedule{}

	err := json.NewDecoder(r.Body).Decode(schedule)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	resp := schedule.Create()
	u.Respond(w, resp)
}

var UpdateStudentAttendance = func(w http.ResponseWriter, r *http.Request) {
	attend := &models.Attendance{}

	err := json.NewDecoder(r.Body).Decode(attend)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	resp := attend.Create()
	u.Respond(w, resp)
}

var GetGroupAttendance = func(w http.ResponseWriter, r *http.Request) {
	month, _ := strconv.Atoi(mux.Vars(r)["month"])
	subId, _ := strconv.Atoi(mux.Vars(r)["sub_id"])
	temp := make([]*models.Schedule, 0)
	err := models.GetDB().Table("schedules").Where("sub_id = ?", subId).Find(&temp).Error
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Error while exec query"))
		return
	}
	l, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Error while load location"))
		return
	}
	date := time.Date(time.Now().Year(), time.Month(month), 1, 0, 0, 0, 0, l)
	resp := u.Message(true, "success")
	for date.Month() == time.Month(month) {
		for i := 0; i < len(temp); i++ {
			if strings.ToLower(temp[i].Day) == strings.ToLower(date.Weekday().String()) {
				t := date.Format("2006-01-02") + " " + temp[i].Time[11:19]
				resp[t] = date.Day()
			}
		}
		date = date.AddDate(0, 0, 1)
	}

	u.Respond(w, resp)
}
