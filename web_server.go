package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Hai struct {
	ID   int    `json:"id"`
	City string `json:"city"`
	Hai  string `json:"hai"`
}

// ------------- Globals -------------
type Hais []Hai

var mainDB *sql.DB

// ------------- Main -------------
func main() {

	db, err := sql.Open("sqlite3", "./database.sqlite")
	mainDB = db

	checkErr(err)

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/hai/", HaiHandler)
	http.HandleFunc("/haicounter/", countHandler)

	fmt.Println("Starting server at port 3333")
	if err := http.ListenAndServe(":3333", nil); err != nil {
		log.Fatal(err)
	}

}

// ------------- Handlers -------------
func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request received")
	io.WriteString(w, "Welcome!")
}

func countHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	rows, err := mainDB.Query("SELECT COUNT(hai), city FROM hais GROUP BY city;")
	checkErr(err)

	type City struct {
		City  string `json:"city"`
		Count int    `json:"count"`
	}

	var cities []City

	for rows.Next() {

		var city City

		err = rows.Scan(&city.Count, &city.City)
		checkErr(err)

		cities = append(cities, city)

	}

	result := struct {
		Result []City `json:"result"`
	}{
		Result: cities,
	}

	json.NewEncoder(w).Encode(result)

}

func HaiHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	m := r.URL.Query()
	query := ""

	id := m.Get("id")
	city := m.Get("city")
	hai := m.Get("hai")

	switch {
	case id != "":
		query = fmt.Sprintf("SELECT * FROM hais where %s == h_id", id)

	case hai != "" && city != "":
		query = "SELECT * FROM hais where hai LIKE '%" + hai + "%' AND city LIKE '%" + city + "%'"

	case hai != "":
		query = "SELECT * FROM hais where hai LIKE '%" + hai + "%'"

	case city != "":
		query = "SELECT * FROM hais where city LIKE '%" + city + "%'"

	default:
		query = fmt.Sprintf("SELECT * FROM hais")
	}

	value := queryDB(query)

	result := struct {
		Result []Hai `json:"result"`
	}{
		Result: value,
	}

	json.NewEncoder(w).Encode(result)

}

// ------------- Database query -------------
func queryDB(query string) []Hai {

	rows, err := mainDB.Query(query)
	checkErr(err)

	var hais Hais
	for rows.Next() {

		var hai Hai

		err = rows.Scan(&hai.ID, &hai.City, &hai.Hai)
		checkErr(err)

		hais = append(hais, hai)

	}

	return hais
}

// ------------- error handler -------------
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
