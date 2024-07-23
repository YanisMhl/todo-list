package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "perntodo"
)

type Todo struct {
	TodoID      string `json:"todo_id"`
	Description string `json:"description"`
}

func initDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getTodos(db *sql.DB) ([]Todo, error) {
	var todos []Todo
	var query string = "SELECT * FROM todo;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.Description, &todo.TodoID)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func getTodo(db *sql.DB, id string) (Todo, error) {
	var todo Todo
	var query string = "SELECT * FROM todo WHERE todo_id=$1;"
	err := db.QueryRow(query, id).Scan(&todo.Description, &todo.TodoID)
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func postTodo(db *sql.DB, description string) error {
	var query string = "INSERT INTO todo (description) VALUES ($1);"
	result, err := db.Exec(query, description)
	fmt.Println("result : ", result)
	return err
}

func putTodo(db *sql.DB, id string, description string) error {
	var query string = "UPDATE todo SET description=$1 WHERE todo_id=$2;"
	result, err := db.Exec(query, description, id)
	fmt.Println("result : ", result)
	return err
}

func deleteTodo(db *sql.DB, id string) error {
	var query string = "DELETE FROM todo WHERE todo_id=$1;"
	result, err := db.Exec(query, id)
	fmt.Println("result : ", result)
	return err
}
