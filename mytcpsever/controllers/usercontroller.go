package controllers

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

//0 成功  1 无此用户  2 密码错误 3 数据库链接错误
//用户登陆成功 返回 用户名 头像地址
func UserLogin(userid string , password string ,lastsessionid string) ResponseData {
	var result ResponseData

	md5Password := enCode(password)
	//判断redis 0 验证成功   1 验证失败  2 redis验证失败（可能第一次登录,redis中没有）

	redisCode := verPassRedis(lastsessionid,md5Password,userid)

	//redis验证失败 去数据库验证
	if redisCode == 2 {
		//从数据库获取数据
		code,nickname,filename := verUserPassword(userid,md5Password)

		//登陆成功，设置sessionid
		if code == 0 {
			result.Nickname = nickname
			result.Filename = filename

			////获取唯一sessionid
			//sessionid := getsessionid()


			//测试专用sessionid
			sessionid := "test"+userid


			var mysession Session
			mysession.Userid = userid
			mysession.Nickname = nickname
			mysession.Filename = filename
			mysession.Password = md5Password
			//根据sessionid将用户的session放入redis中
			SetSession(sessionid,mysession)

			result.Sessionid = sessionid
		}
		result.Resultcode = code

	}else{
		result.Resultcode = redisCode
		result.Sessionid = lastsessionid
	}

	return result
}


//0注册成功 1 用户已经存在 2 系统错误
func UserRegister (userid string, password string,nickname string) int {
	md5Password := enCode(password)
	return userRegisterSql(userid,md5Password,nickname)
}



// 0 获取成功  1 获取失败
func VerUserView (sessionid string,userid string) (int,string,string) {

	verCode,nickname,filename := verUserSession(sessionid,userid)
	if verCode == 0 {
		return 0,nickname,filename
	}else{
		//redis 验证未成功 进入sql验证
		return getViewSql(userid)
	}
}




func getsessionid() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}


//md5单向加密
func enCode(data  string) string {
	h:=md5.New()
	h.Write([]byte (data) )
	return hex.EncodeToString(h.Sum(nil))
}
