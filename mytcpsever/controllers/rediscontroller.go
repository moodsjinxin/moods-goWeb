package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var mypool *redis.Pool

func init(){
	mypool = newpool()
}


func SetSession(sessionid string,session Session){
	sessionJson ,err:= json.Marshal(session)
	if err != nil {
		fmt.Println(err)
	}

	c := mypool.Get()
	defer c.Close()
	_ ,seterr := c.Do("SET",sessionid,sessionJson,"EX","10000")

	if seterr != nil {
		fmt.Print("set redis error!",seterr)
	}
	c.Close()
}

func GetSession(sessionid string) Session{

	c := mypool.Get()
	defer c.Close()
	getsession,geterro := redis.Bytes(c.Do("GET",sessionid))
	if geterro != nil {
		fmt.Print("get redis error!",geterro)
	}
	c.Close()
	var session Session
	err := json.Unmarshal(getsession,&session)
	if err != nil {
		fmt.Println(err)
	}
	return session
}



//0 验证成功   1 验证失败  2 redis验证失败
func verPassRedis(sessionid string,password string,userid string) int {
	c := mypool.Get()
	defer c.Close()
	getsession,geterro := redis.Bytes(c.Do("GET",sessionid))
	if geterro != nil {
		c.Close()
		// 第一次登录 redis 验证失败
		fmt.Println("1处返回 2",geterro)
		return 2
	}

	c.Close()

	var session Session
	err := json.Unmarshal(getsession,&session)
	if err != nil {
		return 2
	}

	// 新用户登录时可能用到了以前用户的cookie传来的session此时判定为新用户登录，交给sql去判断
	if session.Userid != userid {
		return 2
	}

	if session.Password == password {
		return 0
	}
	return 1

}

// 0 获取成功  1 获取失败
func verUserSession(sessionid string,userid string) (int,string,string)  {
	//获取session
	c := mypool.Get()
	defer c.Close()
	getSessionMess,geterro := redis.Bytes(c.Do("GET",sessionid))
	if geterro != nil {
		c.Close()
		//session验证失败，交由sql去执行
		return 1,"",""
	}
	c.Close()

	var session Session
	err := json.Unmarshal(getSessionMess,&session)
	if err != nil {
		return 1,"",""
	}


	if session.Userid != userid {
		return 1,"",""
	}


	return 0,session.Nickname,session.Filename

}



