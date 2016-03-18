package main

import (
	//"database/sql"
	//"flag"
	"fmt"
	//_ "github.com/mattn/go-sqlite3"
	//"log"
	//"os"
)

//func dbHandle(dbFilename string) *sql.DB {
//	db, err := sql.Open("sqlite3", dbFilename)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return db
//}

func main() {
	//db := dbHandle("./games.db")
	//defer db.Close()
	games := Scrape("http://data.ncaa.com/carmen/brackets/championships/basketball-men/d1/2016/data.json")
	fmt.Println(games)
}
