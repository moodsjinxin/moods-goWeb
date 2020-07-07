package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)


//前端数据{userid: , password: , nickname: }
//返回结果 code：0注册成功 1 用户已经存在 2 系统错误
func UserRegister( response http.ResponseWriter, request *http.Request){

	var result Regisresult

	if request.Method == "GET" {
		t,err := template.ParseFiles("WebPages/register.html")
		if err != nil {
			fmt.Println("loading is err!")
		}
		err = t.Execute(response,nil)

	}else if request.Method == "POST" {
		t,err := template.ParseFiles("WebPages/register.html")
		if err != nil {
			fmt.Println("loading is err!")
		}
		// 业务逻辑交由 tcp去做
		var transData Transportdata
		transData.Method = "regist"
		transData.Userid = request.FormValue("userid")
		transData.Password = request.FormValue("password")
		transData.Nickname = request.FormValue("nickname")



		registResponseData := connectTcp(transData)
		result.Code = registResponseData.Resultcode
		if result.Code == 0{
			result.Mess = "注册成功"
		}else{
			result.Mess = "注册失败"
		}

		t.Execute(response,result)

	}else {
		fmt.Println("method is err!")
	}

}


// 请求数据：{userid,password}
// 登录结果数据 code：0 成功  1 失败
// 新登陆用户的cookies：sessionid，userid
func Userlogin( response http.ResponseWriter,request *http.Request) {
	//get时加载页面
	if request.Method == "GET" {
		t,err := template.ParseFiles("WebPages/login.html")
		if err != nil {
			fmt.Println("loading is err!",err)
		}
		err = t.Execute(response,nil)

	}else if request.Method == "POST" {

		Mycookie,_ := request.Cookie("mysession")


		userid := request.FormValue("userid")
		password := request.FormValue("password")


		var transData Transportdata
		transData.Method = "login"
		transData.Userid = userid
		transData.Password = password

		//重复登录增加redis缓存所需要的前端sessionid，新登录用户设置为no
		if Mycookie != nil {
			//将cookie中的值按照分隔符划分为对应字段，（分隔符的选取在set cookie时设置）
			values := strings.Split(Mycookie.Value,"%jinxin%")
			// 借用filename来传输 sessionid
			transData.Filename = values[0]
		}else{
			transData.Filename = "no"
		}


		//调用tcp服务器，获得返回结果
		var responsedata = connectTcp(transData)


		//登陆成功 设置cookie
		if responsedata.Resultcode == 0 {


			//业务逻辑
			//%jinxin% 为分隔符
			cookievalue := responsedata.Sessionid+"%jinxin%"+userid
			//设置cookie
			cookie := http.Cookie{Name: "mysession", Value: cookievalue, Path: "/", MaxAge: 6400}
			http.SetCookie(response, &cookie)

			http.Redirect(response,request,"/view",http.StatusFound)



			////登陆接口测试代码
			//var result LoginResult
			//result.Code = responsedata.Resultcode
			//resultJson, _ := json.Marshal(result)
			//response.Write(resultJson)


		} else {										//登陆失败
			var result LoginResult
			result.Code = responsedata.Resultcode
			resultJson, _ := json.Marshal(result)
			response.Write(resultJson)
		}

	}else {
		fmt.Println("method is err!")
	}
}



//0 验证成功，否则验证失败
func verUser(sessionid string , userid string)  ResponseData {

	var transData Transportdata
	transData.Method = "veruser"
	transData.Userid = userid
	// 借用password字段来传输sessionid
	transData.Password = sessionid

	return connectTcp(transData)

}




