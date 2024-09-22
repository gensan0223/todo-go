package main

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"strings"
	"todo/db"
)

var todos = []Todo{
	{ID: 1, Title: "Task 1", Completed: false},
	{ID: 2, Title: "Task 2", Completed: true},
}

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	db.InitDB()
	mux := http.NewServeMux()

	mux.HandleFunc("/todos", TodosHandler)
	mux.HandleFunc("/todos/", TodoHandler)

	log.Println("Server is running on port 8080!")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func TodosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
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
	case "POST":
		var todo Todo
		_ = json.NewDecoder(r.Body).Decode(&todo)
		todos = append(todos, todo)
		w.Header().Set("Content-type", "application/json")
	}
}

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "GET":
		GetTodo(w, r, id)
	case "PUT":
		UpdateTodo(w, r, id)
	case "DELETE":
		DeleteTodo(w, r, id)
	}
}

func GetTodo(w http.ResponseWriter, r *http.Request, id int) {
	for _, todo := range todos {
		if todo.ID == id {
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(todo)
			return
		}
	}
	http.NotFound(w, r)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request, id int) {
	for i, todo := range todos {
		if todo.ID == id {
			var updatedTodo Todo
			_ = json.NewDecoder(r.Body).Decode(&updatedTodo)
			updatedTodo.ID = id
			todos[i] = updatedTodo
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(todo)
			return
		}
	}
	http.NotFound(w, r)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request, id int) {
	for i, todo := range todos {
		if todo.ID == id {
			todos[i] = todos[len(todos)-1]
			todos = todos[:len(todos)-1]
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}
