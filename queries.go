package main

import (
	"database/sql"
	"fmt"
	"time"
)

//Add a new genre to the database
func insertGenre(db *sql.DB, genre string) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO genres VALUES ( ? )")
	defer stmt.Close()
	checkErr(err)
	res, err := stmt.Exec(genre)
	checkErr(err)
	return res.LastInsertId()
}

//List all genres in the database, in html table
func showAllGenres(db *sql.DB) (string, error) {
	var genre string
	outputTable := "<table>\n<tr><th>Genres</th></tr>\n"
	rows, err := db.Query("SELECT * FROM genres")
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&genre)
		checkErr(err)
		outputTable += "<tr><td>" + genre + "</td></tr>\n"
	}
	outputTable += "</table>\n"
	return outputTable, nil
}

//Add new member
func insertMember(db *sql.DB, surname, firstName, membershipClass string) (int64, error) {

	stmt, err := db.Prepare("INSERT INTO members (surname, first_name, membership_class, membership_started) VALUES ( ?,?,?,? )")
	defer stmt.Close()
	checkErr(err)
	if len(surname) <= 1 || len(firstName) <= 1 {
		return 0, fmt.Errorf("Invalid Name, must contain at least one character")
	}
	classCheck, err := db.Query("SELECT count(*) FROM membership_classes WHERE class = ?", membershipClass)
	defer classCheck.Close()
	//if classcheck == 0 {return 0, fmt.Errorf("Unknown membership class")}
	startDate := time.Now().Format("2006-01-02")
	res, err := stmt.Exec(surname, firstName, membershipClass, startDate)
	checkErr(err)
	return res.LastInsertId()
}
