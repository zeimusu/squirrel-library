package main

import (
	"fmt"
	"net/http"
	"strconv"
)

//Insert a row into a table
func (db *myDB) handleInsert(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err.Error())
		errorHandler(w, r, 500, err.Error())
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
	errorHandler(w, r, 404, "")
}

func insertNewGenre(w http.ResponseWriter, r *http.Request, db *myDB) {
	genreName := r.FormValue("genre")
	if len(genreName) <= 1 {
		errorHandler(w, r, 500, "Genre name empty")
		return
	}
	id, err := insertGenre(db.db, genreName)
	if err != nil {
		errorHandler(w, r, 500, err.Error())
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

		errorHandler(w, r, 500, fmt.Sprintf("invalid membername. got firstname='%v' surname='%v'", firstName, surname))
		return
	}
	//get all membershipclasses from dbase
	//check membershipclass is among them

	id, err := insertMember(db.db, surname, firstName, membershipClass)
	if err != nil {
		errorHandler(w, r, 500, err.Error())
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
		errorHandler(w, r, 500, err.Error())
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
		errorHandler(w, r, 500, err.Error())
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
