package controllers

import (
	"bd_admin/models"
	u "bd_admin/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
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

type subId struct {
	Name  string `json:"sub_name"`
	Group string `json:"group"`
	Id    uint   `json:"sub_id"`
}

var GetSub = func(w http.ResponseWriter, r *http.Request) {
	subs := make([]*models.Schedule, 0)
	user := r.Context().Value("user").(uint)
	err := models.GetDB().Table("schedules").Select("sub_name, sub_id, schedules.group").Where("teacher_id = ?", user).Find(&subs).Error
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Error while exec query"))
		return
	}
	resp := u.Message(true, "success")
	subs2 := make([]*subId, 0)
	for i := 0; i < len(subs); i++ {
		flag := true
		for j := 0; j < len(subs2); j++ {
			if subs[i].SubID == subs2[j].Id {
				flag = false
			}
		}
		if flag {
			temp := subId{}
			temp.Id = subs[i].SubID
			temp.Group = subs[i].Group
			temp.Name = subs[i].SubName
			subs2 = append(subs2, &temp)
		}
	}
	resp["data"] = subs2
	u.Respond(w, resp)
}

type stud struct {
	FirstName string `json:"first_name"`
	Attend    bool   `json:"attend"`
	Date      string `json:"date"`
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
	students := make([]*models.Student, 0)
	k := 0
	for date.Month() == time.Month(month) {
		for i := 0; i < len(temp); i++ {
			if strings.ToLower(temp[i].Day) == strings.ToLower(date.Weekday().String()) {
				err := models.GetDB().Table("students").Where("students.group = ?", temp[i].Group).Find(&students).Error

				if err != nil {
					fmt.Println(err)
					u.Respond(w, u.Message(false, "Error while exec query"))
					return
				}
				t := date.Format("2006-01-02") + " " + temp[i].Time[11:19]
				s := make([]*stud, 0)
				for j := 0; j < len(students); j++ {
					ss := &stud{}
					temp1 := &models.Attendance{}
					err := models.GetDB().Table("attendances").Where("student_id = ? AND sub_id = ? AND date = ?", students[j].AccountID, subId, date.Format("2006-01-02")).First(temp1).Error
					fmt.Println(temp1.StudentId, students[j].AccountID, temp1.Date)

					if err == gorm.ErrRecordNotFound {
						ss.Attend = false
					} else {
						ss.Attend = temp1.Attend
					}
					ss.FirstName = students[j].FirstName
					ss.Date = t
					s = append(s, ss)
				}
				resp[strconv.Itoa(k)] = s
				k++
			}

		}
		date = date.AddDate(0, 0, 1)

	}
	resp["amount"] = k
	u.Respond(w, resp)
}
