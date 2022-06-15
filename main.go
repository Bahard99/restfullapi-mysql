package main

import (
	// "fmt"
	"log"
	// "net/http"
	"encoding/json"
	// "io/ioutil"
	"database/sql"

	// "github.com/gorilla/mux"
	"github.com/gofiber/fiber"
	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	ID 		string	`json:"id"`
	Name	string	`json:"name"`
	Age		int		`json:"age"`
	Grade	int		`json:"grade"`
}

func conn() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_golang")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main()  {
	log.Println("Starting server at localhost:8080")

	app := fiber.New()
	app.Get("/get/:id?", getPerson)
	app.Post("/create", createPerson)
	app.Put("/update/:id", updatePerson)
	app.Delete("/delete/:id", deletePerson)

	app.Listen(8080)

}

func createPerson(c *fiber.Ctx)  {
	db, err := conn()
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	defer db.Close()
	
	var person Person
	json.Unmarshal([]byte(c.Body()), &person)

	_, err = db.Exec("insert into student values (?, ?, ?, ?)", person.ID, person.Name, person.Age, person.Grade)
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	res := "Insert Successfullly!"

	c.Send(res)
}

func getPerson(c *fiber.Ctx) {
	db, err := conn()
	if err != nil {
		c.Status(500).Send(err)
		return 
	}
	defer db.Close()

	query := "select * from student"
	id 	  := ""

	if c.Params("id") != "" {
		id = c.Params("id")
		query = "select * from student where id = ?"
		// query = fmt.Sprintf("%q, %s", query1, id)
		// db.Query("select * from student where id = ?", id)
	}

	rows, err := db.Query(query, id)
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	defer rows.Close()

	var result []Person

	for rows.Next() {
        var each = Person{}
        var err = rows.Scan(&each.ID, &each.Name, &each.Age, &each.Grade)

        if err != nil {
            c.Status(500).Send(err)
		return
        }

        result = append(result, each)
    }

    if err = rows.Err(); err != nil {
        c.Status(500).Send(err)
		return
    }
	if result == nil {
		c.Status(404).Send("Data Not Found")
		return
	}

	json, _ := json.Marshal(result)
	c.Send(json)
}

func updatePerson(c *fiber.Ctx)  {
	db, err := conn()
	if err != nil {
		c.Status(500).Send(err)
		return 
	}
	defer db.Close()

	var person Person
	json.Unmarshal([]byte(c.Body()), &person)

	id := c.Params("id")

	_, err = db.Exec("update student set name = ?, age = ?, grade = ? where id = ?",
						person.Name, person.Age, person.Grade, id)
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	res := "Update Successfully!"

	c.Send(res)
}

func deletePerson(c *fiber.Ctx)  {
	db, err := conn()
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	defer db.Close()

	id := c.Params("id")

	_, err = db.Exec("delete from student where id = ?", id)
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	res := "Delete Successfully!"
	
	c.Send(res)
}