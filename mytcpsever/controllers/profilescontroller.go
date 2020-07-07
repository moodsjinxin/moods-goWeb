package controllers

import (
	"fmt"
	"net"
	"os"
)



//0 上传成功  1 上传失败  2 用户验证失败
//lastdata 是粘包问题导致的残留文件数据
func UploadPicture(data Transportdata,lastDataLen int,lastdata []byte,conn  net.Conn ) int{
	//验证用户
	mysession := GetSession(data.Password)
	if data.Userid != mysession.Userid {
		return 1
	}

	//用户上传的头像放到http服务器的userpictures目录下
	dirpath := "/Users/xinjin/gitproject/entrytask/myhttpsever/userpictures"
	os.MkdirAll(dirpath,os.ModePerm)

	//用当前文件名作图片名
	filepath := dirpath+"/"+ data.Filename
	_,err := os.Stat(filepath)
	if err == nil {
		os.Remove(filepath)
	}

	//创建文件
	file,fileerro := os.Create(filepath)
	if fileerro !=nil{
		fmt.Print("file create erro ：",fileerro)
		return 1
	}

	buf := make([]byte ,data.Len+1024)

	n,_  := conn.Read(buf)

	//当本次获取的数据与传输的文件数据相同时 没有发生粘包
	if n == data.Len {
		_ ,err =file.Write(buf[:n])
	}else{
		_ ,err =file.Write(lastdata[lastDataLen+1:])
		_ ,err =file.Write(buf[:n])
	}
	if err != nil {
		fmt.Println(err)
		return 1
	}

		//先写redis
		lastFileName := mysession.Filename
		mysession.Filename = data.Filename
		SetSession( data.Password, mysession )


		setresult :=setPicture(data.Userid,data.Filename)
		//如果数据库上传头像路径失败，判断本次用户上传失败
		if setresult == 1 {
			//数据库失败，redis还原
			mysession.Filename = lastFileName
			SetSession( data.Password , mysession)
			//删除文件
			os.Remove(filepath)
			return 1
		}

	return 0
}

//修改昵称 0 修改成功  1  修改失败  2 验证失败
func RestNickname(userid string,sessionid string,newnickname string) int {

	mysession := GetSession(sessionid)
	if userid != mysession.Userid {
		return 2
	}

	lastName := mysession.Nickname
	mysession.Nickname = newnickname

	//先修改redis
	SetSession( sessionid, mysession )

	//setcode为设置数据库时返回的结果
	 setcode := Setnickname(userid,newnickname)

	//数据库更改nickname失败(返回1)  -> 用户昵称修改失败
		if setcode == 1{
			//还原redis
			mysession.Nickname = lastName
			SetSession(sessionid,mysession)
			fmt.Println("数据库修改昵称失败")
			return 1
		}else{
			return 0
		}
}

