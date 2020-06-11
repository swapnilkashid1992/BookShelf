package main

import (
	"net/http"

	"./service"
	"github.com/gorilla/mux"
)

/*book1 := models.Book{
	BookName:    "MyBook",
	Auther_Name: "Swapnil",
	IsAvailable: true,
}
book3 := models.Book{
	BookName:    "MyBook3",
	Auther_Name: "Swapnil3",
	IsAvailable: false,
}
book := []models.Book{book1, book3}
service.AddBook(book)
//book = service.ReadAllBook()
//fmt.Println(book)
	fmt.Println(service.FindBookById(2))*/
func main() {
	/*	r := mux.NewRouter()
		r.HandleFunc("/book", service.AddBookTest)
		http.ListenAndServe(":8080", r)
	*/
	r := mux.NewRouter()
	r.HandleFunc("/book", service.AddBook).Methods("POST")
	r.HandleFunc("/book", service.ReadAllBook).Methods("GET")
	r.HandleFunc("/book", service.UpdateBook).Methods("PUT")
	r.HandleFunc("/book/{id}", service.DeleteBook).Methods("DELETE")
	r.HandleFunc("/book/{id}", service.FindBookById).Methods("GET")
	r.HandleFunc("/login", service.Login).Methods("POST")
	r.HandleFunc("/user", service.RegisterUser).Methods("POST")

	http.ListenAndServe(":8081", r)
}
