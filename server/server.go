package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func todosHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getTodosHandler(db, w)
	case "POST":
		postTodoHandler(db, w, r)
	default:
		http.NotFound(w, r)
	}
}

func todoHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getTodoHandler(db, w, r)
	case "PUT":
		putTodoHandler(db, w, r)
	case "DELETE":
		deleteTodoHandler(db, w, r)
	default:
		http.NotFound(w, r)
	}
}

func main() {
	db, err := initDB()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		todosHandler(db, w, r)
	})
	mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		todoHandler(db, w, r)
	})

	fmt.Println("listening on port 8080..")
	err = http.ListenAndServe(":8080", enableCORS(mux))
	if err != nil {
		panic(err)
	}

}
