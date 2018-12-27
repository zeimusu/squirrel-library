package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

//Type to hold the reference to the database. I can put methods on
//this struct to pass it to the various handler functions
type myDB struct {
	db *sql.DB
}

//In Lieu of proper error handling
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)

	}
}

func main() {
	db, err := sql.Open("mysql", "jk:chamame@/library")
	checkErr(err)
	defer db.Close()

	err = db.Ping()
	checkErr(err)

	err = insertBook(db, "9780857503626", "Past Tense", "", "Paperback", "Thriller,Mystery", [][3]string{[3]string{"Child", "Lee", "Author"}}, 1)
	checkErr(err)

	/*
		mydb := &myDB{db: db}
		mydb.server()
		fmt.Println(isbn.To13("0340897301"))
	*/
}
