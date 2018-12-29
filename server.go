package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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

//Show a single table from the database
func (db *myDB) handleShowTable(w http.ResponseWriter, r *http.Request) {
	table := r.URL.Query().Get("table")
	htmlTemplate := `<!doctype html>
	<html>
	<head><title>%s</title></head>
	<body>
	%s
	</body>
	</html>
	`
	if table == "genres" {
		allGenres, err := showAllGenres(db.db)
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(fmt.Sprintf(htmlTemplate, "Genres", allGenres)))
	} else {
		errorHandler(w, r, 404)
	}
}

//Insert a row into a table
func (db *myDB) handleInsert(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err.Error())
		errorHandler(w, r, 500)
		return
	}
	tableName := r.FormValue("table")
	if tableName == "genre" {
		insertNewGenre(w, r, db)
		return
	}
	if tableName == "member" {
		insertNewMember(w, r, db)
		return
	}
	errorHandler(w, r, 404)
}

func insertNewGenre(w http.ResponseWriter, r *http.Request, db *myDB) {
	genreName := r.FormValue("genre")
	if len(genreName) <= 1 {
		errorHandler(w, r, 500)
		return
	}
	id, err := insertGenre(db.db, genreName)
	if err != nil {
		errorHandler(w, r, 500)
	}
	html := `<!doctype html>
<html><head><title>Success</title><head>
<body>
<p>You updated the genre table with %v.</p>
<p>The server returned %v.</p>
<p>Return <a href="/">home</a></p>
</body>
</html>`
	fmt.Fprintf(w, html, genreName, id)
}

func insertNewMember(w http.ResponseWriter, r *http.Request, db *myDB) {

	firstName := r.FormValue("firstName")
	surname := r.FormValue("surname")
	membershipClass := r.FormValue("membershipClass")

	if len(firstName) <= 0 || len(surname) <= 1 {
		fmt.Printf("invalid membername. got firstname='%v' surname='%v'", firstName, surname)
		errorHandler(w, r, 500)
		return
	}
	//get all membershipclasses from dbase
	//check membershipclass is among them

	id, err := insertMember(db.db, surname, firstName, membershipClass)
	if err != nil {
		errorHandler(w, r, 500)
	}
	html := `<!doctype html>
<html><head><title>Success</title><head>
<body>
<p>You updated the members table with %v %v in class %v.</p>
<p>The server returned %v.</p>
<p>Return <a href="/">home</a></p>
</body>
</html>`
	fmt.Fprintf(w, html, firstName, surname, membershipClass, id)
}

func (db *myDB) handleInsertBook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		errorHandler(w, r, 500)
	}
	htmlHead := `<!doctype html>
<html><head><title>Success</title><head>
<body>`
	htmlFoot := "</body></html>"
	w.Write([]byte(htmlHead))
	w.Write([]byte("<table>"))
	for key, value := range r.Form {
		fmt.Fprintf(w, "<tr><td>%v</td><td>%v</td></tr>", key, value)
	}
	w.Write([]byte("</table>"))
	w.Write([]byte(htmlFoot))

	authors := make([][3]string, 0)
	authors = append(authors, [3]string{
		r.FormValue("surname"),
		r.FormValue("firstName"),
		r.FormValue("role"),
	})
	copies, err := strconv.Atoi(r.FormValue("numcopies"))
	if err != nil {
		errorHandler(w, r, 500)
	}

	err = insertBook(
		db.db,
		r.FormValue("ISBN"),
		r.FormValue("title"),
		r.FormValue("dewey"),
		r.FormValue("format"),
		r.FormValue("genrelist"),
		authors,
		copies,
	)
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
	http.HandleFunc("/insertbook", db.handleInsertBook)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
