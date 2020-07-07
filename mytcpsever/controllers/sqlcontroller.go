package controllers

import (
	"database/sql"
	"fmt"
	_"github.com/go-sql-driver/mysql"
	"time"
)

//var db,dbconerr = sql.Open("mysql","root:jx32380078@tcp(127.0.0.1:3306)/moods?charset=utf8")

var db *sql.DB
func init(){
	db,_ = sql.Open("mysql","root:jx32380078@tcp(127.0.0.1:3306)/moods?charset=utf8")
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(30)
	db.SetConnMaxLifetime(time.Second*180)
	db.Ping()
}


//0 成功  1 无此用户  2 密码错误 3 数据库链接错误
// 返回验证结果及nickname 若失败则nickname为""
func verUserPassword(userid string,mypassword string) (int,string,string){
	//db,dbconerr := sql.Open("mysql","root:jx32380078@tcp(127.0.0.1:3306)/moods?charset=utf8")
	rows,qerro := db.Query("select password,nickname,filename from user_info where userid = ?", userid )
	defer rows.Close()
	if qerro !=nil {
		fmt.Print(qerro)
		return 1,"err","err"
	}
	var getpassword string
	var getnickname string
	var getfilename string
	for rows.Next(){
		err := rows.Scan(&getpassword,&getnickname,&getfilename)
		if err !=nil {
			fmt.Print(err)
			rows.Close()
			return 3,"err","err"
		}
	}

	if getpassword == mypassword{
		return 0,getnickname,getfilename
	}
	rows.Close()
	return 2,"",""
}


//0注册成功 1 用户已经存在 2 系统错误
func userRegisterSql(userid string, password string, nickname string) int {

	//db,dbconerr := sql.Open("mysql","root:jx32380078@tcp(127.0.0.1:3306)/moods?charset=utf8")，已设置为全局

	rows,qerro := db.Query("select userid from user_info where userid = ?", userid )
	defer rows.Close()
	if qerro != nil {
		fmt.Print(qerro)
		return 2
	}

	//判断查询的用户是否存在
	for rows.Next(){
		var result string
		err := rows.Scan(&result)
		if err != nil{
			fmt.Print(err)
		}
		// 用户名已经存在，（注意，此处需要前端匹配 用户名不能为 nil）
		if result != "nil" {
			rows.Close()
			return 1
		}
	}

	// 用户不存在，插入新的用户数据，新用户的头像为默认"default"，设置为默认头像
	rows,err :=db.Query("insert into user_info values (?,?,?,?)",userid,password,nickname,"default")
	if err != nil {
		fmt.Print(err)
	}
	rows.Close()
	return 0
}


// 1 数据库中修改图像路径失败 0 更改成功
func setPicture(userid string,filepath  string) int {

	rows,err := db.Query("update user_info set filename = ? where  userid = ?",filepath,userid)
	defer rows.Close()
	if err != nil {
		rows.Close()
		return 1
	}
	rows.Close()
	return 0
}

//修改昵称 0 修改成功  1  修改失败
func Setnickname(userid string, newnicknae string) int {

	rows,err := db.Query("update user_info set nickname = ? where  userid = ?",newnicknae,userid)
	if err != nil {
		rows.Close()
		return 1
	}
	rows.Close()
	return 0
}


//从sql中获取个人信息
func getViewSql(userid string) (int,string,string){
	rows,qerro := db.Query("select nickname,filename from user_info where userid = ?", userid )
	defer rows.Close()
	if qerro !=nil {
		fmt.Print(qerro)
		return 1,"",""
	}

	var getnickname string
	var getfilename string
	for rows.Next(){
		err := rows.Scan(&getnickname,&getfilename)
		if err !=nil {
			rows.Close()
			return 1,"",""
		}
	}
	rows.Close()
	return 0,getnickname,getfilename
}
