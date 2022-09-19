package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{
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

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func getEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookID := mux.Vars(r)["id"]

	num := 0 // this number helps me handle not found cases
	for _, book := range events {
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

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookID := mux.Vars(r)["id"]

	num := 0 // this number helps me handle not found cases
	for index, book := range events {
		if book.ID == bookID {
			events = append(events[:index], events[index+1:]...)
			num++
		}
	}
	// this block of code handles not found cases
	if num == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "book not found")
	}
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Kindly enter data with event title and description only in order to update")
	}

	var updateBook event
	bookID := mux.Vars(r)["id"]
	json.Unmarshal(reqBody, &updateBook)
	num := 0 //this number helps me handle not found cases
	for index, book := range events {
		if book.ID == bookID {
			theBook := &events[index]
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
	router.HandleFunc("/books", createEvent).Methods("POST")
	router.HandleFunc("/books", getAllEvents).Methods("GET")
	router.HandleFunc("/books/{id}", getEvent).Methods("GET")
	router.HandleFunc("/books/{id}", deleteEvent).Methods("DELETE")
	router.HandleFunc("/books/{id}", updateEvent).Methods("PATCH")

	log.Fatal(http.ListenAndServe(":8000", router))
}
