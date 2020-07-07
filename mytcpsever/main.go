package main

import (
	"./controllers"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

var wg =sync.WaitGroup{}

func main(){
	listencer,err := net.Listen("tcp","127.0.0.1:8001")
	if err != nil{
		fmt.Print("tcp connect erro !",err)
		return
	}

	for{

		conn,err := listencer.Accept()
		if err != nil {
			fmt.Print("accept erro :",err)
			continue
		}
		go requestexec(conn)

	}

}

func requestexec(conn net.Conn){
	var data controllers.Transportdata
	var responseData controllers.ResponseData
	jsonMessage := make([]byte, 1024)


	//defer conn.Close()
	for {
		//tcp 获取jsaon数据方式1
		//jsonMessage := make([]byte, 1024)
		_, err := conn.Read(jsonMessage)
		if err != nil {
			continue
		}

		//byte[0] 表示请求的数据长度，用来截取粘包情况下的图片数据
		datalen := jsonMessage[0]

		//if err != nil {
		//	fmt.Println(err)
		//}
		//var data controllers.Transportdata
		err = json.Unmarshal(jsonMessage[1:datalen+1], &data)
		if err != nil {
			fmt.Println(err)
		}


		//// tcp 获取json数据方式2，此方法不知道怎么解决粘包问题
		//d  := json.NewDecoder(conn)
		//var data controllers.Transportdata
		//fmt.Println(data)
		//var result ResponseData
		//err := d.Decode(&data)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//fmt.Println(data.Method)

		//var responseData controllers.ResponseData
		if data.Method == "login" {                  //调用登录功能
			responseData = controllers.UserLogin(data.Userid, data.Password,data.Filename)

			responDataJson,err := json.Marshal(responseData)
			if err != nil {
				fmt.Println(err)
			}
			conn.Write(responDataJson)
			//break
		} else if data.Method == "regist" {          //调用注册方法

			responseData.Resultcode = controllers.UserRegister(data.Userid,data.Password,data.Nickname)

			responDataJson,err := json.Marshal(responseData)
			if err != nil {
				fmt.Println(err)
			}
			conn.Write(responDataJson)
			//break
		}else if data.Method == "upload" {          //调用文件上传功能
			responseData.Resultcode = controllers.UploadPicture(data,int(datalen),jsonMessage,conn)

			responDataJson,err := json.Marshal(responseData)
			if err != nil {
				fmt.Println(err)
			}
			conn.Write(responDataJson)
			//break
		}else if data.Method == "resetnickname" {

			responseData.Resultcode = controllers.RestNickname(data.Userid,data.Password,data.Nickname)

			responDataJson,err := json.Marshal(responseData)
			if err != nil {
				fmt.Println(err)
			}
			conn.Write(responDataJson)
			//break
		}else if  data.Method == "veruser" {
			//借用的password传送sessionid
			responseData.Resultcode,responseData.Nickname,responseData.Filename = controllers.VerUserView(data.Password,data.Userid)

			responDataJson,err := json.Marshal(responseData)
			if err != nil {
				fmt.Println(err)
			}
			conn.Write(responDataJson)
			//break
		}
		//break
	}
}

