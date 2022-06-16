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
	ID 		string	`json:"id,omitempty"`
	Name	string	`json:"name,omitempty"`
	Age		int		`json:"age,omitempty"`
	Grade	int		`json:"grade,omitempty"`
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
	app.Get("/getall", getAllPerson)
	app.Get("/getone/:id", getOnePerson)
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

func getOnePerson(c *fiber.Ctx) {
	db, err := conn()
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	defer db.Close()

	id := c.Params("id")

	var each = Person{}

	err = db.QueryRow("select * from student where idd = ?", id).Scan(&each.ID, &each.Name, &each.Age, &each.Grade)
	switch {
	case err == sql.ErrNoRows:
		c.Status(404).Send("Data Not Found!!")
		return
	case err != nil:
		c.Status(500).Send(err)
		return
	}

	json, _ := json.Marshal(each)
	c.Send(json)
}
 
func getAllPerson(c *fiber.Ctx) {
	db, err := conn()
	if err != nil {
		c.Status(500).Send(err)
		return 
	}
	defer db.Close()

	rows, err := db.Query("select * from student")
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