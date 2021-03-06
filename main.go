package main

import (
	"bd_admin/app"
	"bd_admin/controllers"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/student/info", controllers.CreateStudent).Methods("POST")
	router.HandleFunc("/api/teacher/info", controllers.CreateTeacher).Methods("POST")
	router.HandleFunc("/api/sub/new", controllers.CreateSub).Methods("POST")
	router.HandleFunc("/api/teacher/setAttendance", controllers.UpdateStudentAttendance).Methods("POST")
	router.HandleFunc("/api/teacher/qrcode", controllers.GetQRCode).Methods("POST")

	router.HandleFunc("/api/teacher/sub", controllers.GetSub).Methods("GET")
	router.HandleFunc("/api/student/qrcode", controllers.SentQRCode).Methods("GET")
	router.HandleFunc("/api/student/getAllAttendance", controllers.GetStudentAttendance).Methods("GET")
	router.HandleFunc("/api/user/getInfo", controllers.GetInfo).Methods("GET")
	router.HandleFunc("/api/teacher/{month}/{sub_id}", controllers.GetGroupAttendance).Methods("GET")

	router.Use(app.JwtAuthentication) //attach JWT auth middleware
	router.Use(app.TeacherRights)
	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
