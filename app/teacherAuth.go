package app

import (
	"bd_admin/models"
	u "bd_admin/utils"
	"github.com/jinzhu/gorm"
	"net/http"
)

var TeacherRights = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/user/new", "/api/user/login", "/api/student/info", "/api/student/getAllAttendance", "/api/student/qrcode"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                                                                                                            //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		user := r.Context().Value("user").(uint)
		//acc := &models.Account{}
		temp := &models.Account{}

		response := make(map[string]interface{})
		err := models.GetDB().Table("accounts").Where("ID = ?", user).First(temp).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			response = u.Message(false, "Connection Error")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		if temp.Role == "student" {
			response = u.Message(false, "account hasn't such permission")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		//Everything went well
		next.ServeHTTP(w, r)
	})
}
