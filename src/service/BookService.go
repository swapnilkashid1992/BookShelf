package service

import (
	"encoding/json"
	"log"
	"net/http"

	"../models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func CreateDbConnection() *gorm.DB {
	db, err := gorm.Open("postgres", "user=postgres dbname=BOOKSHELF sslmode=disable password=admin")
	if err != nil {
		log.Panic(err)
	}
	return db
}
func CloseConnection(db *gorm.DB) {
	db.Close()
}
func AddBook(w http.ResponseWriter, r *http.Request) {
	db := CreateDbConnection()
	defer CloseConnection(db)
	if !db.HasTable(&models.Book{}) {
		db.CreateTable(&models.Book{})
	}
	var books []models.Book
	json.NewDecoder(r.Body).Decode(&books)
	for _, v := range books {
		var presentBook models.Book
		db.Where("book_name=?", v.BookName).Find(&presentBook)
		if (models.Book{}) == presentBook {
			db.Create(&v)
		} else {
			log.Println("Book Already Present")
		}
	}
}
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	db := CreateDbConnection()
	defer CloseConnection(db)
	var book models.Book
	json.NewDecoder(r.Body).Decode(&book)
	var presentBook models.Book
	db.Where("book_name=?", book.BookName).Find(&presentBook)
	if (models.Book{}) != presentBook {
		db.Update(&book)
	} else {
		log.Println("Book Already Present")
	}
}
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	db := CreateDbConnection()
	defer CloseConnection(db)
	vars := mux.Vars(r)
	id := vars["id"]
	var presentBook models.Book
	db.Where("id=?", id).Find(&presentBook)
	if (models.Book{}) != presentBook {
		db.Delete(&presentBook)
	} else {
		log.Println("Book Already Present")
	}
}

func ReadAllBook(w http.ResponseWriter, r *http.Request) {
	db := CreateDbConnection()
	defer CloseConnection(db)
	var books []models.Book
	db.Find(&books)
	b, _ := json.Marshal(books)
	w.Write(b)
}

func FindBookById(w http.ResponseWriter, r *http.Request) {
	db := CreateDbConnection()
	defer CloseConnection(db)
	var presentBook models.Book
	vars := mux.Vars(r)
	id := vars["id"]
	db.Where("id=?", id).Find(&presentBook)
	if (models.Book{}) == presentBook {
		log.Println("Book Does Not Present")
	}
	b, _ := json.Marshal(presentBook)
	w.Write(b)
}
