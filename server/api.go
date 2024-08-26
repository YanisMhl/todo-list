package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func getParam(url string) (string, error) {
	var param string = strings.TrimPrefix(url, "/todos/")
	if strings.Contains(param, "/") {
		return "", errors.New("not a valid path")
	}
	return param, nil
}

func extractUserIDFromToken(tokenString string) (string, error) {
	secret := []byte("supersecret")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["id"].(string)
		return userID, nil
	} else {
		return "", errors.New("invalid token")
	}
}

func loginUserHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Name == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	response, err := loginUser(db, user.Name, user.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		fmt.Println(err)
		return
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func signUpUserHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Println("c'est la merde")
		fmt.Println(err)
		return
	}
	user, err = signUpUser(db, user.Name, user.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		fmt.Println("error : ", err)
		return
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
}

func getTodosHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
	fmt.Println(authHeader)
	tokenString := authHeader[len("Bearer : "):]
	userID, err := extractUserIDFromToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	todos, err := getTodos(db, userID)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		fmt.Println("error : ", err)
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
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
	fmt.Println(authHeader)
	tokenString := authHeader[len("Bearer : "):]
	userID, err := extractUserIDFromToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		fmt.Println("erreur : ", err)
		return
	}
	fmt.Println("description : ", todo.Description)
	err = postTodo(db, userID, todo.Description)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	w.Write([]byte("successfully created todo"))
}

func putTodoHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var todo Todo
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
	fmt.Println(authHeader)
	tokenString := authHeader[len("Bearer : "):]
	_, err := extractUserIDFromToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&todo)
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
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
	fmt.Println(authHeader)
	tokenString := authHeader[len("Bearer : "):]
	_, err := extractUserIDFromToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
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
