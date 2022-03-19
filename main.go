package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Connection to database
func getMySQLDS() *sql.DB {
	db, err := sql.Open("mysql", "root@(127.0.0.1:3306)/go?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFiles("crudForm.html"))
}

type studentInfo struct {
	Sid    string
	Name   string
	Course string
}

func studentHandle(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDS()
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	student := studentInfo{
		Sid:    r.FormValue("sid"),
		Name:   r.FormValue("name"),
		Course: r.FormValue("course"),
	}

	if r.FormValue("clicked") == "Insert" {
		sid, _ := strconv.Atoi(student.Sid)
		_, err := db.Exec("insert into studentinfo(sid, name, course) values(?, ?, ?)", sid, student.Name, student.Course)

		if err != nil {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: err.Error()})
		} else {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: "Record Inserted"})

		}

	} else if r.FormValue("clicked") == "Read" {
		data := []string{}
		rows, err := db.Query("select * from studentinfo")

		if err != nil {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: err.Error()})
		} else {
			s := studentInfo{}
			data = append(data, "<table border=2>")
			data = append(data, "<tr><th>Student ID</th><th>Name</th><th>Course</th></tr>")
			for rows.Next() {
				rows.Scan(&s.Sid, &s.Name, &s.Course)
				data = append(data, fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>", s.Sid, s.Name, s.Course))
			}
			data = append(data, "</table>")
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: strings.Trim(fmt.Sprint(data), "[]")})
			// fmt.Println(studata)
		}

	} else if r.FormValue("clicked") == "Update" {
		tmpl.Execute(w, struct {
			Success bool
			Message string
		}{true, "Data Updated"})

	} else if r.FormValue("clicked") == "Delete" {
		tmpl.Execute(w, struct {
			Success bool
			Message string
		}{true, "Data Deleted"})

	}
	// fmt.Println(student)
	// fmt.Println(r.FormValue("clicked"))

}
func main() {

	http.HandleFunc("/", studentHandle)

	// starting the server
	http.ListenAndServe(":3000", nil)
}
