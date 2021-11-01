package controllers

import (
	"bd_admin/models"
	u "bd_admin/utils"
	"encoding/json"
	"net/http"
)

var CreateStudent = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user") . (uint) //Получение идентификатора пользователя, отправившего запрос
	student := &models.Student{}

	err := json.NewDecoder(r.Body).Decode(student)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	student.AccountID = user
	resp := student.Create()
	u.Respond(w, resp)
}