package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (db *myDB) handleBorrow(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		errorHandler(w, r, 500, err.Error())
	}

	borrower := r.FormValue("memberID")
	borrowerID, err := strconv.Atoi(borrower)
	if err != nil {
		errorHandler(w, r, 500, err.Error())
		return
	}
	book := r.FormValue("bookID")
	bookID, err := strconv.Atoi(book)
	if err != nil {
		errorHandler(w, r, 500, err.Error())
		return
	}
	due, err := borrowBook(db.db, int64(bookID), int64(borrowerID))
	if err != nil {
		fmt.Fprintln(w, `<html><!doctype html>
<html><head><title>Failure</title><head>
<body>`)
		fmt.Fprintln(w, "<h1>Borrow failure</h1><p>Server responded:</p>")
		fmt.Fprintf(w, "<blockquote>%v</blockquote></body></html>", err.Error())
		return
	}
	fmt.Fprintln(w, `<html><!doctype html>
<html><head><title>Success</title><head>
<body>`)
	fmt.Fprintf(w, "<h1>Borrow success</h1><p>Due date: %v.</p>", due)
}

func (db *myDB) handleReturn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		errorHandler(w, r, 500, err.Error())
		return
	}

	book := r.FormValue("bookID")
	bookID, err := strconv.Atoi(book)
	if err != nil {
		errorHandler(w, r, 500, err.Error())
		return
	}

	overdueDays, err := returnBook(db.db, int64(bookID))

	if err != nil {
		fmt.Fprintln(w, `<html><!doctype html>
<html><head><title>Failure</title><head>
<body>`)
		fmt.Fprintln(w, "<h1>Borrow failure</h1><p>Server responded:</p>")
		fmt.Fprintf(w, "<blockquote>%v</blockquote></body></html>", err.Error())
		return
	}
	fmt.Fprintln(w, `<html><!doctype html>
<html><head><title>Success</title><head>
<body>`)
	fmt.Fprintln(w, "<h1>Book returned success</h1>")
	if overdueDays > 0 {
		fmt.Fprintf(w, "<p>This book is %v days overdue.</p>", overdueDays)
	} else {
		fmt.Fprintf(w, "<p>This book is not overdue.</p>")
	}
	fmt.Println("</body></html>")

}
