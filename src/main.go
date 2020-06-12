package main

import (
	"net/http"

	"../src/service"
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
	r.HandleFunc("/login", service.Login).Methods("POST")
	r.HandleFunc("/user", service.RegisterUser).Methods("POST")
	s := r.PathPrefix("/api").Subrouter()
	s.Use(service.JwtRegularVerify)
	s.HandleFunc("/book", service.ReadAllBook).Methods("GET")
	s.HandleFunc("/book/{id}", service.FindBookById).Methods("GET")
	s.HandleFunc("/booking/{id}/user/{uid}", service.BookABook).Methods("POST")
	s.HandleFunc("/booking/{id}", service.DeleteBooking).Methods("DELETE")

	p := r.PathPrefix("/api").Subrouter()
	p.Use(service.JwtAdminVerify)
	p.HandleFunc("/book", service.AddBook).Methods("POST")
	p.HandleFunc("/book", service.ReadAllBook).Methods("GET")
	p.HandleFunc("/book", service.UpdateBook).Methods("PUT")
	p.HandleFunc("/book/{id}", service.DeleteBook).Methods("DELETE")
	p.HandleFunc("/book/{id}", service.FindBookById).Methods("GET")
	http.ListenAndServe(":8081", r)

}
