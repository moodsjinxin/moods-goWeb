package controllers


type Regisresult struct {
	Code  int
	Mess string
}

type LoginResult struct {
	Code int
}

//tcp和http服务器交互的通用数据格式
type Transportdata struct {
	Method string
	Userid string
	Nickname string
	Filename string
	Password string
	Len      int
}

type ResponseData struct {
	Resultcode int
	Nickname string
	Filename string
	Sessionid string
	Userid  string
	Mess   string
}
