package main

import (
	"./controllers"
	"log"
	"net/http"
)


func main(){
	//路由、接口
	http.HandleFunc("/setnickname",controllers.Resetnickname)
	http.HandleFunc("/upload",controllers.UploadHandler)
	http.HandleFunc("/register",controllers.UserRegister)
	http.HandleFunc("/login",controllers.Userlogin)
	http.HandleFunc("/view/",controllers.ProfileView)
	http.Handle("/", http.FileServer(http.Dir("userpictures")))
	log.Fatal(http.ListenAndServe("127.0.0.1:8000",nil))
}
