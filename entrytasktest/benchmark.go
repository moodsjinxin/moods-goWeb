package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var taskNum,concurrencyNum,isRandom int

func init(){
	flag.IntVar(&taskNum,"num",4000,"total run task num")  //返回指针类型，绑定到第一个参数的变量上
	flag.IntVar(&concurrencyNum,"co",200,"concurren num")
	flag.IntVar(&isRandom,"random",0,"random:1(y)/0(n)")
}

func main(){
	flag.Parse()
	totalTaskNum := int(taskNum)     //总请求数目
	totalConcurrNum := int(concurrencyNum)         //并发量

	//var transport http.RoundTripper = &http.Transport{
	//	DialContext:(&net.Dialer{
	//		Timeout: 30 * time.Second,
	//		KeepAlive: 30 * time.Second,
	//	}).DialContext,
	//	MaxIdleConns:  totalConcurrNum,
	//	MaxIdleConnsPerHost: totalConcurrNum,
	//	IdleConnTimeout: 90*time.Second,
	//	TLSHandshakeTimeout: 10*time.Second,
	//	ExpectContinueTimeout: 1*time.Second,
	//}

	//client := &http.Client{
	//	Transport: transport,
	//}

	var  runTaskNum = int32(0)
	var wg sync.WaitGroup

	username := "test100"

	for i:=0;i<totalConcurrNum;i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()

			for atomic.AddInt32(&runTaskNum,1) <= int32(totalTaskNum) {

				if isRandom == 1{
					username = "test" + strconv.Itoa(rand.Intn(10000000))
				}

				fmt.Println(username)

				//测试时写死sessionid
				sessionid := "test"+username

				data := url.Values{}

				data.Set("userid",username)
				//测试所用的数据账户和密码相同
				data.Set("password",username)


				//req,_ := http.NewRequest("POST","http://127.0.0.1/login",bytes.NewBufferString(data.Encode()))

				req,_ := http.NewRequest("POST","http://127.0.0.1:8000/login",strings.NewReader(data.Encode()))
				req.AddCookie(&http.Cookie{Name: "mysession",Value: sessionid+"%jinxin%"+username,Expires: time.Now().Add(120*time.Second)})
				req.Header.Add("Content-Type","application/x-www-form-urlencoded")
				req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))


				client01 := &http.Client{}
				resp,_ := client01.Do(req)

				fmt.Println(resp)
				resp.Body.Close()
			}
		}(i)
	}

	timestart := time.Now()
	wg.Wait()

	timeElapsed := time.Since(timestart)
	qps := float64(totalTaskNum)/timeElapsed.Seconds()


	fmt.Printf("Benchmark -  IsRandom: %v   ToalTaskNum: %v   TotalConcurrNum: %v TimeElapsed: %v  QPS: %v",
		isRandom,totalTaskNum,totalConcurrNum,timeElapsed,math.Ceil(qps))

}
