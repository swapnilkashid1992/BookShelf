package service

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"sync"

	"../models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var wg sync.WaitGroup

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
	err := r.ParseMultipartForm(5 * 1024 * 1024)
	if err != nil {
		panic(err)
	}
	fmt.Println("Hello")
	file, handler, err := r.FormFile("fileupload")

	if file == nil {
		json.NewDecoder(r.Body).Decode(&books)
	} else {
		books = parseCsv(file, handler.Filename)
	}
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
func parseCsv(file multipart.File, fileName string) []models.Book {
	f, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)

	io.Copy(f, file)
	if err != nil {
		log.Println("Error in opening file")
	}
	f.Close()
	f, err = os.Open("C:\\Users\\gs-1454\\Desktop\\GolagWorkspace\\BookShelf\\src\\booklist.csv")
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	records, _ := csv.NewReader(f).ReadAll()
	var books []models.Book
	for _, row := range records {
		//	for _, detail := range row {
		//	details := strings.Split(detail, ",")
		availability, _ := strconv.ParseBool(row[2])
		book := models.Book{
			BookName:    row[0],
			Auther_Name: row[1],
			IsAvailable: availability,
		}
		books = append(books, book)
	}

	return books
}
