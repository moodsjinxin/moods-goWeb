package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)


//cookie 存放的是 sessionid，userid
func ProfileView(w http.ResponseWriter, r  *http.Request){

	Mycookie,err := r.Cookie("mysession")

	//将cookie中的值按照分隔符划分为对应字段，（分隔符的选取在set cookie时设置）
	values := strings.Split(Mycookie.Value,"%jinxin%")

	if err != nil{
		http.Error(w, " ProfileView get cookies err", http.StatusInternalServerError)
		return
	}else{
		//从session中获取用户的nickname、filename并展示到页面上
		verData := verUser(values[0],values[1])
		if verData.Resultcode == 0 {
			t,err := template.ParseFiles("WebPages/profiles.html")
			if err != nil {
				http.Error(w, " ProfileView loading is err!"+err.Error(), http.StatusInternalServerError)
				return
			}
			if verData.Filename == "default"{
				verData.Filename = "/default.jpg"
			}else {
				verData.Filename = "/"+verData.Filename
			}
			verData.Userid = values[1]
			verData.Mess = r.URL.Path[len("/view"):]
			err = t.Execute(w,verData)

		}else if verData.Resultcode == 1 {
			http.Error(w, " ProfileView 用户验证失败", http.StatusInternalServerError)
			return
		}
	}
}

//cookie 存放的是 sessionid，userid
func UploadHandler(w http.ResponseWriter, r  *http.Request){

	Mycookie,err := r.Cookie("mysession")
	values := strings.Split(Mycookie.Value,"%jinxin%")
	if err != nil{
		http.Error(w, " 认证失败，请重新登录", http.StatusInternalServerError)
		return
	}else{

		picture,head,_ := r.FormFile("filename")

		if picture == nil{
			http.Redirect(w,r,"/view/nopicture",http.StatusFound)
		}else{
			picturename := head.Filename

			var transData Transportdata
			transData.Method = "upload"
			transData.Userid = values[1]
			transData.Filename = picturename

			//借用password字段来传sessionid
			transData.Password = values[0]
			code := sendFileToTcp(transData,picture,int(head.Size)).Resultcode

			if code == 0 {
				http.Redirect(w,r,"/view/success",http.StatusFound)
			}else{
				http.Redirect(w,r,"/view/failed",http.StatusFound)
			}
		}
	}
}


//修改昵称 返回的code为0 修改成功  1  修改失败  2  sessionid验证失败
func Resetnickname(w http.ResponseWriter, r  *http.Request){

	Mycookie,err := r.Cookie("mysession")
	values := strings.Split(Mycookie.Value,"%jinxin%")
	if err != nil{
		fmt.Fprint(w,"认证失败，请重新登陆")
		return
	}else{
		newnickname := r.FormValue("newnickname")
		if newnickname == ""{
			http.Redirect(w,r,"/view/name-empty",http.StatusFound)
		}else{
			var transData Transportdata
			transData.Method = "resetnickname"
			transData.Userid = values[1]
			transData.Nickname = newnickname
			transData.Password = values[0]  // sessionid

			resultCode := connectTcp(transData).Resultcode

			if resultCode == 0 {
				http.Redirect(w,r,"/view/success",http.StatusFound)
			}else{
				http.Redirect(w,r,"/view/failed",http.StatusFound)
			}
		}

	}
}
