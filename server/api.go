package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func getParam(url string) (string, error) {
	var param string = strings.TrimPrefix(url, "/todos/")
	if strings.Contains(param, "/") {
		return "", errors.New("not a valid path")
	}
	return param, nil
}

func getTodosHandler(db *sql.DB, w http.ResponseWriter) {
	todos, err := getTodos(db)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	todosJSON, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(todosJSON)
}

func getTodoHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	param, err := getParam(r.URL.String())
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	todo, err := getTodo(db, param)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	todoJSON, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(todoJSON)
}

func postTodoHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		fmt.Println("erreur : ", err)
		return
	}
	fmt.Println("description : ", todo.Description)
	err = postTodo(db, todo.Description)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	w.Write([]byte("successfully created todo"))
}

func putTodoHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		fmt.Println("erreur : ", err)
		return
	}
	id, err := getParam(r.URL.String())
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	err = putTodo(db, id, todo.Description)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	w.Write([]byte("successfully updated todo"))
}

func deleteTodoHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	id, err := getParam(r.URL.String())
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	err = deleteTodo(db, id)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	w.Write([]byte("successfully deleted todo"))
}
