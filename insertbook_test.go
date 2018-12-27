package main

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func compareslices(a, b []string) bool {
	if len(a) != len(b) {
		fmt.Println(len(a), len(b))
		return false
	}

	for i, g := range a {
		if !(b[i] == g) {
			fmt.Printf("*%v*  *%v*", b[i], g)
			return false
		}
	}
	return true
}
func TestValidateGenres(t *testing.T) {
	db, err := sql.Open("mysql", "jk:chamame@/library")
	checkErr(err)
	defer db.Close()

	validGenres := []string{"Horror", "Child fiction", "Western"}
	validWrongCase := []string{"HORROR", "comedy"}
	invalidGenres := []string{"Horror", "Rhubarb"}

	err = validateGenres(db, &validGenres)
	if err != nil {
		t.Errorf("Valid genres found invalid")
	}
	if !compareslices(validGenres,
		[]string{"Horror", "Child fiction", "Western"}) {
		t.Errorf("valid genres changed to %v, from %v", validGenres, []string{"Horror", "Child fiction", "Western"})
	}

	err = validateGenres(db, &validWrongCase)
	if err != nil {
		t.Errorf("Valid but wrong case genres found invalid")
	}
	if !compareslices(validWrongCase, []string{"Horror", "Comedy"}) {
		t.Errorf("wrong genres is %v", validWrongCase)
	}

	err = validateGenres(db, &invalidGenres)
	if err == nil {
		t.Errorf("invalid genres return no error")
	}
}
