package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net"
	"time"
)


var cp  *ConnPool

func init() {
	cp,_ = NewConnPool(func() (ConnRes, error) {
		return net.Dial("tcp", "127.0.0.1:8001");
	}, 2000, time.Second*180);  //1600
	fmt.Println(cp.len())
}


//所有的请求都调用tcp服务器去完成
func connectTcp(transData Transportdata) ResponseData{
	var responsedata ResponseData

	transDataJson,err := json.Marshal(transData)
	if err != nil {
		fmt.Println(err)
	}

	//封装请求数据包，bytes【0】存放的是数据长度,解决文件传输的粘包问题
	datalen := len(transDataJson)
	var requestDatas  []byte
	//加入数据长度
	requestDatas = append(requestDatas, byte(datalen))
	//数组合并
	combinebytes := [][]byte{requestDatas,transDataJson}
	requestDatas = bytes.Join(combinebytes,[]byte{})


	////与服务端建立链接,非链接池方法
	//serverAdrr,err := net.ResolveTCPAddr("tcp","127.0.0.1:8001")
	//conn,err := net.DialTCP("tcp",nil,serverAdrr)
	//if err != nil {
	//	fmt.Println("链接tcp服务器失败：",err)
	//}

	connres,err := cp.Get()
	//实例化
	conn := connres.(*net.TCPConn)


	_, err = conn.Write(requestDatas)
	if err != nil {
		fmt.Println(err)
	}

	// 获取服务器的响应
	d := json.NewDecoder(conn)
	err = d.Decode(&responsedata)
	if err != nil {
		fmt.Println(err)
	}
	cp.Put(connres)

	return responsedata
}



func sendFileToTcp(transData Transportdata,picture multipart.File,maxlen int ) ResponseData{
	var responsedata ResponseData

	buf := make([]byte ,maxlen+1024)
	filelen,err := picture.Read(buf)
	if err != nil {
		fmt.Println(err)
	}

	transData.Len = filelen

	transDataJson,err := json.Marshal(transData)
	if err != nil {
		fmt.Println(err)
	}

	// 封装数据包 bytes【0】存放的是数据长度
	datalen := len(transDataJson)
	var requestDatas  []byte
	//加入数据长度
	requestDatas = append(requestDatas, byte(datalen))
	//数组合并
	combinebytes := [][]byte{requestDatas,transDataJson}
	requestDatas = bytes.Join(combinebytes,[]byte{})

	//serverAdrr,err := net.ResolveTCPAddr("tcp","127.0.0.1:8001")
	//conn,err := net.DialTCP("tcp",nil,serverAdrr)
	//if err != nil {
	//	fmt.Println("链接tcp服务器失败：",err)
	//}

	connres,err := cp.Get()
	conn := connres.(*net.TCPConn)
	//传输请求数据
	_, err = conn.Write(requestDatas)
	if err != nil {
		fmt.Println(err)
	}


	//传输文件数据，注意粘包问题，使用的简单的添加数据长度确定文件数据的位置
	_,err = conn.Write(buf[:filelen])
	if err != nil {
		fmt.Println(err)
	}


	// 获取服务器的响应
	d := json.NewDecoder(conn)

	//放回链接池
	cp.Put(connres)

	err = d.Decode(&responsedata)
	if err != nil {
		fmt.Println(err)
	}


	return responsedata

}


