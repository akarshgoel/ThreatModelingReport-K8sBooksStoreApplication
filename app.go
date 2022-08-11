package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// book struct
type Book struct {
	ID        string  `json:"id"`
	Publisher string  `json:"publisher"`
	Title     string  `json:"title"`
	Author    *Author `json:"author"`
}

// author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// init books variable as slice book struct
var books []Book

//get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

//get book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	value := mux.Vars(r) // get book id
	// loop through books and find with book id
	for _, item := range books {
		if item.ID == value["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

//create book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(1000000)) // test id
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

//delete book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	value := mux.Vars(r)
	for index, item := range books {
		if item.ID == value["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

//update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	value := mux.Vars(r)
	for index, item := range books {
		if item.ID == value["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = value["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	// init router
	r := mux.NewRouter()

	// test data
	books = append(books, Book{ID: "01", Publisher: "Test1 Publications", Title: "Book 01", Author: &Author{Firstname: "Test1", Lastname: "User1"}})
	books = append(books, Book{ID: "02", Publisher: "Test2 Publications", Title: "Book 02", Author: &Author{Firstname: "Test2", Lastname: "User2"}})
	books = append(books, Book{ID: "03", Publisher: "Test3 Publications", Title: "Book 03", Author: &Author{Firstname: "Test3", Lastname: "User3"}})

	// route handlers for endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServeTLS(":9000", "localhost.crt", "localhost.key", r))
}

