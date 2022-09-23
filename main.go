package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type book struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type allbooks []book

// this acts as the application database
var books = allbooks{
	{
		ID:          "1",
		Title:       "Introduction to golang",
		Description: "Come join us for a chance to see how golang works and eventually get to try it",
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Println("i am running")
	fmt.Fprintf(w, "welcome home!")
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var newBook book
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newBook)
	books = append(books, newBook)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newBook)
}

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookID := mux.Vars(r)["id"]

	num := 0 // this number helps me handle not found cases
	for _, book := range books {
		if book.ID == bookID {
			json.NewEncoder(w).Encode(book)
			num++
		}
	}

	// this block of code handles not found cases
	if num == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "book not found")
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookID := mux.Vars(r)["id"]

	num := 0 // this number helps me handle not found cases
	for index, book := range books {
		if book.ID == bookID {
			books = append(books[:index], books[index+1:]...)
			fmt.Fprintln(w, "book deleted")
			num++
		}
	}
	// this block of code handles not found cases
	if num == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "book not found")
	}
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Kindly enter data with event title and description only in order to update")
	}

	var updateBook book
	bookID := mux.Vars(r)["id"]
	json.Unmarshal(reqBody, &updateBook)
	num := 0 //this number helps me handle not found cases
	for index, book := range books {
		if book.ID == bookID {
			theBook := &books[index]
			theBook.Title = updateBook.Title
			theBook.Description = updateBook.Description
			json.NewEncoder(w).Encode(theBook)
			num++
		}
	}
	// this block of code handles not found cases
	if num == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "book not found")
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books", getAllBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	router.HandleFunc("/books/{id}", updateBook).Methods("PATCH")

	log.Fatal(http.ListenAndServe(":8000", router))
}
