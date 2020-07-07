package controllers


type ResponseData struct {
	Resultcode int
	Nickname string
	Filename string
	Sessionid string
	Userid string
	Mess  string
}

type Transportdata struct {
	Method string
	Userid string
	Nickname string
	Filename string
	Password string
	Len int
}

type Session struct {
	Nickname string
	Userid  string
	Filename string
	Password string
}

type ViewData struct {
	Nickname string
	Filename string
}