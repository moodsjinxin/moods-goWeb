package controllers

import (
	"github.com/garyburd/redigo/redis"
	"time"
)



func newConnectReids()(redis.Conn ,error){
	return redis.Dial("tcp","127.0.0.1:6379")
}


func newpool()(*redis.Pool){
	return &redis.Pool{
		MaxIdle: 35,   //pool的最大链接
		Dial: newConnectReids, //链接redis
		MaxActive: 35,
		Wait: true,  //当链接超出后，是否等待其他链接释放
		IdleTimeout: 180*time.Second,
	}
}