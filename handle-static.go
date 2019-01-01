package main

import (
	"fmt"
	"net/http"
)

//Default page handler
func handleLibrary(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

//Error page handler
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "custom 404")
	}
}

func handleNewGenre(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "newgenre.html")
}

func handleNewMember(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "newmember.html")
}

func handleNewBook(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "newbook.html")
}

func handleBorrowReturn(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "borrow-return.html")
}
