package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "perntodo"
)

type User struct {
	UserID    string    `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type Todo struct {
	TodoID      string    `json:"todo_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      string    `json:"user_id"`
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

func loginUser(db *sql.DB, username string, password string) (map[string]string, error) {
	var user User
	var userQuery string = "SELECT id, password FROM users WHERE name=$1;"
	err := db.QueryRow(userQuery, username).Scan(&user.UserID, &user.Password)
	if err != nil {
		return nil, errors.New("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.UserID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte("supersecret"))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response := map[string]string{"token": tokenString, "name": username}
	return response, nil
}

func signUpUser(db *sql.DB, username string, password string) (User, error) {
	var nameResult string
	nameQuery := "SELECT name FROM users WHERE name=$1;"
	nameRow := db.QueryRow(nameQuery, username)
	err := nameRow.Scan(&nameResult)
	if err == nil {
		return User{}, errors.New("this user already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return User{}, err
	}
	var query string = "INSERT INTO users (name, password) VALUES ($1, $2);"
	_, err = db.Exec(query, username, string(hashedPassword))
	if err != nil {
		return User{}, err
	}
	var user User
	var userQuery string = "SELECT * FROM users WHERE name=$1;"
	result := db.QueryRow(userQuery, username).Scan(&user.UserID, &user.Name, &user.Password, &user.CreatedAt)
	fmt.Println(result)
	return user, nil
}

func getTodos(db *sql.DB, userID string) ([]Todo, error) {
	var todos []Todo
	var query string = "SELECT * FROM todo WHERE user_id=$1 ORDER BY created_at;"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.Description, &todo.TodoID, &todo.CreatedAt, &todo.UserID)
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
	err := db.QueryRow(query, id).Scan(&todo.Description, &todo.TodoID, &todo.CreatedAt)
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func postTodo(db *sql.DB, userID string, description string) error {
	var query string = "INSERT INTO todo (user_id, description) VALUES ($1, $2);"
	result, err := db.Exec(query, userID, description)
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
