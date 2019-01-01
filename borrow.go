package main

import (
	"database/sql"
	"fmt"
	"time"
)

func borrowBook(db *sql.DB, bookID, memberID int64) (string, error) {
	var memberClass, surname, firstName string
	var loanLength int64
	err := db.QueryRow(`
	SELECT surname, first_name, membership_class, loan_length     
	FROM members JOIN membership_classes         
	ON members.membership_class=membership_classes.class 
	WHERE id=?`,
		memberID).Scan(&surname, &firstName, &memberClass, &loanLength)
	if err != nil {
		return "", err
	}
	//check for banned user
	//calculate due date
	now := time.Now()
	due := now.AddDate(0, 0, int(loanLength))
	var otherBorrower sql.NullInt64
	var otherDue sql.NullString
	var isbn, title string
	err = db.QueryRow(`
	SELECT member_id, due, books.ISBN, title
	FROM books JOIN book_info ON books.ISBN = book_info.ISBN
	WHERE id=?`, bookID).Scan(&otherBorrower, &otherDue, &isbn, &title)
	if err != nil {
		return "", err
	}
	if otherBorrower.Valid {
		return "", fmt.Errorf("Book %v (%v) is already on loan to %v, until %v", title, bookID, otherBorrower.Int64, otherDue.String)
	}
	//check book is kosher
	stmt, err := db.Prepare(`
	UPDATE books
	SET member_id = ?, due = ?
	WHERE id = ?`)
	defer stmt.Close()
	if err != nil {
		return "", err
	}
	_, err = stmt.Exec(memberID, due, bookID)
	if err != nil {
		return "", err
	}
	return due.Format("2006-01-02"), nil
}

func returnBook(db *sql.DB, bookID int64) (int, error) {

	var borrower sql.NullInt64
	var due sql.NullString
	var isbn, title string
	err := db.QueryRow(`
	SELECT member_id, due, books.ISBN, title
	FROM books JOIN book_info ON books.ISBN = book_info.ISBN
	WHERE id=?`, bookID).Scan(&borrower, &due, &isbn, &title)
	if err != nil {
		return 0, err
	}
	if !borrower.Valid {
		return 0, fmt.Errorf("This book (%v) does not appear to be on loan", title)
	}
	dueDate, err := time.Parse("2006-01-02", due.String)
	if err != nil {
		return 0, err
	}
	overdue := time.Since(dueDate)
	overdueDays := int(overdue.Hours() / 24)
	if overdueDays < 0 {
		overdueDays = 0
	}
	stmt, err := db.Prepare(`
	UPDATE books
	SET member_id = NULL, due = NULL
	WHERE id = ?`)
	defer stmt.Close()
	if err != nil {
		return 0, err
	}
	_, err = stmt.Exec(bookID)
	if err != nil {
		return 0, err
	}
	return overdueDays, nil
}
