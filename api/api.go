package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/shawlyahsan/api_server/auth"

	"github.com/gorilla/mux"
)

//Book Struct

type Book struct {
	ID    string `json:"id"`
	ISBN  string `json:"isbn"`
	Title string `json:"title"`
}

var books []Book
var book_id int
var mutex sync.Mutex

func login(w http.ResponseWriter, r *http.Request) {
	log.Println("Login")
	w.Write([]byte("Login Successful!\n"))

	token_string, err := auth.Get_Token()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could Not Generate A Token" + err.Error()))
	} else {
		w.Header().Set("Authorization", "Bearer "+token_string)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Token: " + token_string))
	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, book := range books {
		if book.ID == params["id"] {
			json.NewEncoder(w).Encode(book)
			return
		}
	}

	json.NewEncoder(w).Encode(&Book{})
}

func createBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(book_id)
	book_id++

	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for idx, book := range books {
		if book.ID == params["id"] {
			var new_book Book
			_ = json.NewDecoder(r.Body).Decode(&new_book)
			new_book.ID = params["id"]
			books[idx] = new_book
			break
		}
	}

	json.NewEncoder(w).Encode(books)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for idx, book := range books {
		if book.ID == params["id"] {
			books = append(books[:idx], books[idx+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(books)
}

func HandleRequestFunc() {

	//init router
	r := mux.NewRouter()

	//sample books
	books = append(books, Book{ID: "1", ISBN: "100000", Title: "Book 1"})
	books = append(books, Book{ID: "2", ISBN: "200000", Title: "Book 2"})
	books = append(books, Book{ID: "3", ISBN: "300000", Title: "Book 3"})
	book_id = 4

	//route handlers / endpoints
	r.HandleFunc("/login", auth.BasicAuthentication(login)).Methods("POST")
	r.HandleFunc("/books", auth.Is_Authorized(getBooks)).Methods("GET")
	r.HandleFunc("/books/{id}", auth.Is_Authorized(getBook)).Methods("GET")
	r.HandleFunc("/books", auth.Is_Authorized(createBook)).Methods("POST")
	r.HandleFunc("/books/{id}", auth.Is_Authorized(updateBook)).Methods("PUT")
	r.HandleFunc("/books/{id}", auth.Is_Authorized(deleteBook)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
