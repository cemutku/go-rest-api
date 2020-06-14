package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/cemutku/go-rest-api/models"
	"github.com/gorilla/mux"
)

func handleHome(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Welcome to the Todo API")
}

// TodoController todo controller
func handleTodos(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case "GET":
		getTodos(rw, req)
	case "POST":
		postTodo(rw, req)
	default:
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("not found"))
	}
}

func handleTodo(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case "GET":
		getTodo(rw, req)
	case "PUT":
		putTodo(rw, req)
	case "DELETE":
		deleteTodo(rw, req)
	default:
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("not found"))
	}
}

// getTodos get all todos
func getTodos(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	encodeResponse(models.GetTodos(), rw)
}

// getTodo get single todo
func getTodo(rw http.ResponseWriter, req *http.Request) {
	reqParams := mux.Vars(req)

	if val, ok := reqParams["id"]; ok {
		id, err := strconv.Atoi(val)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		todo, err := models.GetTodoByID(id)

		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		encodeResponse(todo, rw)
	}
}

// postTodo create a todo
func postTodo(rw http.ResponseWriter, req *http.Request) {
	todo, err := parseRequest(req)

	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		rw.Write([]byte("Unable to parse object!"))
		return
	}

	newTodo, err := models.AddTodo(todo)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	encodeResponse(newTodo, rw)
}

// putTodo update a todo
func putTodo(rw http.ResponseWriter, req *http.Request) {
	reqParams := mux.Vars(req)

	if val, ok := reqParams["id"]; ok {
		id, err := strconv.Atoi(val)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		todo, err := parseRequest(req)

		if err != nil {
			rw.WriteHeader(http.StatusUnprocessableEntity)
			rw.Write([]byte("Unable to parse object!"))
			return
		}

		if todo.ID != id {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Id of todo must match id in url"))
			return
		}

		updatedTodo, err := models.UpdateTodo(todo)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		encodeResponse(updatedTodo, rw)
	}
}

// deleteTodo delete a todo
func deleteTodo(rw http.ResponseWriter, req *http.Request) {
	reqParams := mux.Vars(req)

	if val, ok := reqParams["id"]; ok {
		id, err := strconv.Atoi(val)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		error := models.DeleteTodo(id)

		if error != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	}
}

func encodeResponse(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

func parseRequest(req *http.Request) (models.Todo, error) {
	decoder := json.NewDecoder(req.Body)
	var todo models.Todo
	err := decoder.Decode(&todo)
	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

// RegisterTodoControllers controller route
func RegisterTodoControllers() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handleHome)
	router.HandleFunc("/todos", handleTodos)
	router.HandleFunc("/todos/{id}", handleTodo).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8082", router))
}
