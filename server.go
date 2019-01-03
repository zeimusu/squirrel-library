package main

import (
	"fmt"
	"log"
	"net/http"
)

func (db *myDB) server() {
	fmt.Println("Library server")
	/*
		Create a webserver on port 8080
	*/
	http.HandleFunc("/", handleLibrary)
	http.HandleFunc("/showtable", db.handleShowTable)
	http.HandleFunc("/insert", db.handleInsert)
	http.HandleFunc("/newgenre", handleNewGenre)
	http.HandleFunc("/newmember", handleNewMember)
	http.HandleFunc("/newbook", handleNewBook)
	http.HandleFunc("/borrow-return", handleBorrowReturn)
	http.HandleFunc("/borrowbook", db.handleBorrow)
	http.HandleFunc("/returnbook", db.handleReturn)
	http.HandleFunc("/insertbook", db.handleInsertBook)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
