package main


import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)


var db *sql.DB
func init(){
	db, _ = sql.Open("mysql","root:jx32380078@tcp(127.0.0.1:3306)/moods?charset=utf8")
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(70)
	db.Ping()
}


func SqlInserte(){
	//insert(0,1000000)
	// insert(1000000,5000000)
	insert(5000000,10000000)
}


func insert(start int,end int ){
	//var db,dbconerr = sql.Open("mysql","root:jx32380078@tcp(127.0.0.1:3306)/moods?charset=utf8")

	fmt.Println(start)

	for i:=start;i<end ;i=i+5{
		rows,err := db.Query("insert into user_info values(?,?,?,?),(?,?,?,?),(?,?,?,?),(?,?,?,?),(?,?,?,?)",
			"test"+strconv.Itoa(i),"test"+strconv.Itoa(i),"testname","default","test"+strconv.Itoa(i+1),"test"+strconv.Itoa(i+1),"testname","default","test"+strconv.Itoa(i+2),"test"+strconv.Itoa(i+2),"testname","default","test"+strconv.Itoa(i+3),"test"+strconv.Itoa(i+3),"testname","default","test"+strconv.Itoa(i+4),"test"+strconv.Itoa(i+4),"testname","default")
		defer rows.Close()
		if err!= nil {
			fmt.Println(err)
			break
		}
		rows.Close()
	}
	fmt.Println(end)
}

