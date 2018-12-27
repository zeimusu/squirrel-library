package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/moraes/isbn"
)

//Split a list of genres by commas
func parseGenres(g string) []string {
	return strings.Split(g, ",")
}

//Check each genre is in the db,
//do caseless lookup, then replace with case as in table
//select * from users where upper(first_name) = 'FRED';
func validateGenres(db *sql.DB, g *[]string) error {
	var dbGenre string
	for i, genre := range *g {
		ucgenre := strings.ToUpper(strings.TrimSpace(genre))
		err := db.QueryRow("SELECT genre FROM genres WHERE upper(genre)= ?", ucgenre).Scan(&dbGenre)
		if err != nil {
			return err
		}
		//if //row is empty {return fmt.Errorf("Unknown genre %v",genre)}
		(*g)[i] = dbGenre
	}
	return nil
}

//Split a dewey decimal code such as "510.63" to 510, 63
//check for errors such as non numeric code or more than one decimal point.
func parseDewey(d string) (int64, int64, error) {
	codes := strings.Split(d, ".")
	if len(codes) > 2 {
		return 0, 0, fmt.Errorf("Invalid Dewey code %v.", d)
	}
	major, err := strconv.Atoi(codes[0])
	if err != nil {
		return 0, 0, err
	}
	if len(codes) == 1 {
		return int64(major), 0, nil
	}
	minor, err := strconv.Atoi(codes[1])
	if err != nil {
		return int64(major), 0, err
	}
	return int64(major), int64(minor), nil
}

//Clean an isbn string by removing any spaces hyphens or other characters
//Check it validates as either isbn10 or isbn13
//Return the isbn13 code
func cleanAndCheckIsbn(ISBN string) (string, error) {
	var b strings.Builder
	for _, char := range ISBN {
		if (char >= '0' && char <= '9') || char == 'X' {
			b.WriteRune(char)
		}
	}
	s := b.String()
	if len(s) == 10 && isbn.Validate10(s) {
		isbn13, err := isbn.To13(s)
		if err != nil {
			return "", err
		}
		return isbn13, nil
	}
	if len(s) == 13 && isbn.Validate13(s) {
		return s, nil
	}
	return "", fmt.Errorf("An invalid isbn code %s was passed", s)
}

func checkFormat(f string) (string, error) {
	formats := map[string]string{
		"HARDBACK":  "Hardback",
		"PAPERBACK": "Paperback",
		"CD":        "CD",
	}

	format, ok := formats[strings.TrimSpace(strings.ToUpper(f))]
	if ok {
		return format, nil
	}
	return "", fmt.Errorf("format '%v' not found", f)

}

func getAuthorId(db *sql.DB, a [3]string) int64 {
	surname := strings.TrimSpace(a[0])
	firstName := strings.TrimSpace(a[1])
	var id int64
	err := db.QueryRow("SELECT id FROM authors WHERE surname=? and first_names=?", surname, firstName).Scan(&id)
	if err != nil {
		stmt, err := db.Prepare("INSERT INTO authors (surname, first_names) VALUES (?,?)")
		if err != nil {
			checkErr(err)
		}
		res, err := stmt.Exec(surname, firstName)
		if err != nil {
			checkErr(err)
		}
		newId, err := res.LastInsertId()
		if err != nil {
			checkErr(err)
		}
		return newId
	}
	return id

}

func insertBookInfo(db *sql.DB, ISBN, title string, deweyMajor, deweyMinor int64, format string) error {
	stmt, err := db.Prepare("INSERT INTO book_info (ISBN, title, dewey_major, dewey_minor, format) VALUES (?,?,?,?,?);")
	if err != nil {
		return err
	}
	nullableDeweyMajor := sql.NullInt64{deweyMajor, true}
	nullableDeweyMinor := sql.NullInt64{deweyMinor, true}
	if deweyMajor == 0 {
		nullableDeweyMajor.Valid = false
		nullableDeweyMinor.Valid = false
	}
	_, err = stmt.Exec(ISBN, title, nullableDeweyMajor, nullableDeweyMinor, format)
	return err
}

func insertBookGenres(db *sql.DB, ISBN string, genres []string) error {
	stmt, err := db.Prepare("INSERT INTO book_genre_mapping (ISBN, genre) VALUES (?, ?);")
	if err != nil {
		return err
	}

	for _, g := range genres {
		_, err = stmt.Exec(ISBN, g)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertBookAuthors(db *sql.DB, ISBN string, authorMap map[int64][3]string) error {
	stmt, err := db.Prepare("INSERT INTO book_author_mapping (ISBN, author_id, role) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	for aid := range authorMap {
		_, err = stmt.Exec(ISBN, aid, authorMap[aid][2])
		if err != nil {
			return err
		}
	}
	return nil
}

func insertCopies(db *sql.DB, ISBN string, copies int) error {
	if copies < 1 {
		return fmt.Errorf("Cannot add less than one copy of the book")
	}
	stmt, err := db.Prepare("INSERT INTO books (ISBN) VALUES (?)")
	if err != nil {
		return err
	}
	for i := 0; i < copies; i++ {
		_, err = stmt.Exec(ISBN)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertBook(
	db *sql.DB,
	ISBN,
	title,
	dewey,
	format,
	genres string,
	authors [][3]string,
	copies int) error {
	//validate and clean ISBN
	cleanIsbn, err := cleanAndCheckIsbn(ISBN)
	if err != nil {
		return err
	}
	//check if isbn already in db
	isbn_row := db.QueryRow(
		"SELECT title FROM book_info WHERE ISBN=?", cleanIsbn)
	var existingTitle string
	if err := isbn_row.Scan(&existingTitle); err == nil {
		return fmt.Errorf("The isbn %v is already in the database with title %s %v.", ISBN, existingTitle, cleanIsbn)
	}

	//split genres
	genreList := parseGenres(genres)

	if err := validateGenres(db, &genreList); err != nil {
		return err
	}

	//if a dewey code is given, split it into major and minor parts
	var major, minor int64
	if len(dewey) > 0 {
		major, minor, err = parseDewey(dewey)
		if err != nil {
			return err
		}
	}

	//check that the format is acceptable, correcting case
	format, err = checkFormat(format)
	if err != nil {
		return err
	}

	authorMap := make(map[int64][3]string)
	//check if each author in author table
	for _, a := range authors {
		aid := getAuthorId(db, a)
		authorMap[aid] = a
	}

	err = insertBookInfo(db, ISBN, title, major, minor, format)
	if err != nil {
		return err
	}

	err = insertBookGenres(db, ISBN, genreList)
	if err != nil {
		return err
	}

	err = insertBookAuthors(db, ISBN, authorMap)
	if err != nil {
		return err
	}

	err = insertCopies(db, ISBN, copies)
	if err != nil {
		return err
	}

	return nil
}
