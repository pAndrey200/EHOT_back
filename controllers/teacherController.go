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
	attend2 := &models.Attendance{}
	err := json.NewDecoder(r.Body).Decode(attend)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	err = models.GetDB().Table("attendances").Where("student_id = ? AND sub_id = ? AND date = ?", attend.StudentId, attend.SubId, attend.Date).First(attend2).Error
	var resp map[string]interface{}
	if err == gorm.ErrRecordNotFound {
		resp = attend.Create()
	} else {
		err = models.GetDB().Table("attendances").Where("student_id = ? AND sub_id = ? AND date = ?", attend.StudentId, attend.SubId, attend.Date).Update("attend", attend.Attend).Error
		if err != nil {
			fmt.Println(err)
			u.Respond(w, u.Message(false, "Error while exec query"))
			return
		}
		resp = u.Message(true, "success update")

	}

	u.Respond(w, resp)
}

type str struct {
	Key string `json:"key"`
}

var GetQRCode = func(w http.ResponseWriter, r *http.Request) {
	s := str{}
	err := json.NewDecoder(r.Body).Decode(&s)

	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	words := strings.Fields(s.Key)
	attend := &models.Attendance{}
	attend.StudentId, _ = strconv.ParseUint(words[0], 10, 64)
	attend.Attend, _ = strconv.ParseBool(words[1])
	attend.Date = time.Now().Format("2006-01-02")

	curTime := time.Now()
	s1 := models.Schedule{}
	ss := models.Student{}
	err = models.GetDB().Table("students").Where("account_id = ?", words[0]).First(&ss).Error
	if err != nil {
		u.Respond(w, u.Message(false, "Error while exec query "+err.Error()))
		return
	}
	err = models.GetDB().Table("schedules").Where("schedules.group = ? AND time > ? AND time < ? AND day = ?", ss.Group, curTime.Add(-40*time.Minute+3*time.Hour).Format("15:04:05"), curTime.Add(10*time.Minute+3*time.Hour).Format("15:04:05"), strings.ToLower(time.Now().Weekday().String())).First(&s1).Error
	if err == gorm.ErrRecordNotFound {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "there is no lesson now"))
		return
	} else {
		attend.SubId = uint64(s1.SubID)
	}

	attend2 := &models.Attendance{}
	err = models.GetDB().Table("attendances").Where("student_id = ? AND sub_id = ? AND date = ?", attend.StudentId, attend.SubId, attend.Date).First(attend2).Error
	if err == gorm.ErrRecordNotFound {
		attend.Create()
	}
	resp := u.Message(true, "success")
	u.Respond(w, resp)
}

type SubId struct {
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
	subs2 := make([]*SubId, 0)
	for i := 0; i < len(subs); i++ {
		flag := true
		for j := 0; j < len(subs2); j++ {
			if subs[i].SubID == subs2[j].Id {
				flag = false
			}
		}
		if flag {
			temp := SubId{}
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
	StudentId uint   `json:"student_id"`
	Attend    bool   `json:"attend"`
	Date      string `json:"date"`
}

var ChetnostOfWeek = func(date time.Time) int {
	l, _ := time.LoadLocation("Europe/Moscow")
	startDate := time.Date(time.Now().Year(), time.Month(2), 7, 0, 0, 0, 0, l)
	k := 1
	_, st := startDate.ISOWeek()
	_, en := date.ISOWeek()
	if st > en {
		return -2
	}
	for st != en {
		k++
		st++
	}
	if k%2 == 1 {
		return 1
	} else {
		return -1
	}

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
	startDate := time.Date(time.Now().Year(), time.Month(2), 7, 0, 0, 0, 0, l)
	_, st := startDate.ISOWeek()
	_, en := date.ISOWeek()
	resp := u.Message(true, "success")
	students := make([]*models.Student, 0)
	k := 0
	for date.Month() == time.Month(month) {
		for i := 0; i < len(temp); i++ {
			if (strings.ToLower(temp[i].Day) == strings.ToLower(date.Weekday().String())) && (temp[i].Week == 0 || temp[i].Week == ChetnostOfWeek(date)) && (en >= st) {
				err := models.GetDB().Table("students").Where("students.group = ?", temp[i].Group).Find(&students).Error
				fmt.Println(temp[i].Group)
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
					//fmt.Println(temp1.StudentId, students[j].AccountID, temp1.Date)

					if err == gorm.ErrRecordNotFound {
						ss.Attend = false
					} else {
						ss.Attend = temp1.Attend
					}
					ss.FirstName = students[j].FirstName
					ss.StudentId = students[j].AccountID
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
