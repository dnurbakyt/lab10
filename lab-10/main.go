package main

import (
	"log"

	"database/sql"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type server struct {
	db *sql.DB
}

type Car struct {
	Id        int
	FirstName string
	LastName  string
	CarModel  string
	Price     int
	Hours     int
}

func database() server {
	database, _ := sql.Open("sqlite3", "carsharing.db")
	server := server{db: database}
	return server
}

func (s *server) rent(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		fn := r.FormValue("fn")
		ln := r.FormValue("ln")
		cm := r.FormValue("cm")
		price := r.FormValue("price")
		hours := r.FormValue("hours")

		_, err := s.db.Exec("INSERT INTO carsharing(firstName, lastName, carModel, price, hours) VALUES ($1, $2, $3, $4, $5)", fn, ln, cm, price, hours)

		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	t, _ := template.ParseFiles("static/rent.html")
	t.Execute(w, nil)

}

func (s *server) rented(w http.ResponseWriter, r *http.Request) {
	var Cars []Car
	result, _ := s.db.Query("select * from carsharing;")
	for result.Next() {
		var car Car
		err := result.Scan(&car.Id, &car.FirstName, &car.LastName, &car.CarModel, &car.Price, &car.Hours)
		if err != nil {
			log.Fatal(err)
		}
		Cars = append(Cars, car)
	}
	tmpl, _ := template.ParseFiles("static/rented.html")
	tmpl.Execute(w, Cars)
}

func (s *server) updateCar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		fn := r.FormValue("fn")
		ln := r.FormValue("ln")
		cm := r.FormValue("cm")
		price := r.FormValue("price")
		hours := r.FormValue("hours")
		_, err := s.db.Exec("UPDATE carsharing SET firstName=$1, lastName=$2, carModel=$3, price=$4, hours=$5 WHERE id=$6", fn, ln, cm, price, hours, id)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseFiles("static/update.html")
	t.Execute(w, nil)
}

func (s *server) deleteCar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		_, _ = s.db.Exec("delete from carsharing where id=$1", id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	t, _ := template.ParseFiles("static/delete.html")
	t.Execute(w, nil)
}

func main() {
	s := database()

	fileServer := http.FileServer(http.Dir("./static"))

	http.Handle("/", fileServer)
	http.HandleFunc("/rent", s.rent)
	http.HandleFunc("/rented", s.rented)
	http.HandleFunc("/updateCar", s.updateCar)
	http.HandleFunc("/deleteCar", s.deleteCar)

	defer s.db.Close()
	http.ListenAndServe(":8080", nil)

}
