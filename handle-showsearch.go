package main

import (
	"fmt"
	"log"
	"net/http"
)

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
