package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"todo/db"

	_ "github.com/lib/pq"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	db.InitDB()
	defer db.GetDB().Close()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /todos", GetTodos)
	mux.HandleFunc("POST /todos", PostTodos)
	mux.HandleFunc("GET /todos/{id}", GetTodoById)
	mux.HandleFunc("PUT /todos/{id}", UpdateTodo)
	mux.HandleFunc("DELETE /todos/{id}", DeleteTodo)

	log.Println("Server is running on port 8080!!!")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.GetDB().Query("SELECT id, title, completed FROM todos")
	if err != nil {
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed); err != nil {
			return
		}
		todos = append(todos, todo)
	}
	json.NewEncoder(w).Encode(todos)
}

func PostTodos(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	query := `INSERT INTO todos VALUES($1, $2, $3) RETURNING id, title, completed`
	err := db.GetDB().QueryRow(query, todo.ID, todo.Title, todo.Completed).Scan(&todo.ID, &todo.Title, &todo.Completed)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func GetTodoById(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	id, _ := strconv.Atoi(r.PathValue("id"))
	query := `SELECT * FROM todos WHERE id = $1`
	err := db.GetDB().QueryRow(query, id).Scan(&todo.ID, &todo.Title, &todo.Completed)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("no user with id %d\n", id)
	case err != nil:
		log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	id, _ := strconv.Atoi(r.PathValue("id"))
	query := `UPDATE todos SET title = $1, completed = $2 WHERE id = $3 RETURNING id, title, completed`
	err := db.GetDB().QueryRow(query, todo.Title, todo.Completed, id).Scan(&todo.ID, &todo.Title, &todo.Completed)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("no user with id %d\n", id)
	case err != nil:
		log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	id, _ := strconv.Atoi(r.PathValue("id"))
	query := `DELETE FROM todos WHERE id = $1`
	_, err := db.GetDB().Exec(query, id)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("no user with id %d\n", id)
	case err != nil:
		log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(todo)
}
