package controllers

import (
	"bd_admin/models"
	u "bd_admin/utils"
	"encoding/json"
	"net/http"
)

var CreateTeacher = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user") . (uint) //Получение идентификатора пользователя, отправившего запрос
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
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	resp := attend.Create()
	u.Respond(w, resp)
}