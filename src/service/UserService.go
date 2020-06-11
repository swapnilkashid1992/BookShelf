package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"../models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	jwt.StandardClaims
	Role string `json:"role"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	db := CreateDbConnection()
	defer CloseConnection(db)
	if !db.HasTable(&models.User{}) {
		db.CreateTable(&models.User{})
	}
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user)
	bpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	if err != nil {
		log.Panic(err)
	}
	user.Password = string(bpass)
	user.Role = "regular"
	var u models.User
	db.Where("username=?", user.Username).Find(&u)
	if (models.User{}) != u {
		log.Println("User with this username already Exists")
	} else {
		db.Create(&user)
	}

}

func BookABook(w http.ResponseWriter, r *http.Request) {
	db := CreateDbConnection()
	defer CloseConnection(db)
	if !db.HasTable(&models.Booking{}) {
		db.CreateTable(&models.Booking{})
	}
	vars := mux.Vars(r)
	userId, _ := strconv.ParseInt(vars["uid"], 10, 0)
	bookid, _ := strconv.ParseInt(vars["id"], 10, 0)
	var book models.Book
	db.Where("id = ?", bookid).Find(&book)
	if (models.Book{}) == book {
		log.Println("Book Doesnt present")
		json.NewEncoder(w).Encode("Book Doesnt present")
		return
	}
	if book.IsAvailable == false {
		log.Println("Book is not available")
		json.NewEncoder(w).Encode("Book is not available")
		return
	}
	var b models.Booking
	db.Where("user_id = ? AND deleted_at IS NULL", userId).Find(&b)
	if (models.Booking{}) == b {
		booking := models.Booking{
			USERID: int(userId),
			BOOKID: int(bookid),
		}
		db.Create(&booking)
		json.NewEncoder(w).Encode("Book is done")
	} else {
		log.Println("Book is not available")
		json.NewEncoder(w).Encode("Book is not available")
	}
}

func DeleteBooking(w http.ResponseWriter, r *http.Request) {
	db := CreateDbConnection()
	defer CloseConnection(db)

	vars := mux.Vars(r)
	bookid, _ := strconv.ParseInt(vars["id"], 10, 0)
	var booking models.Booking
	db.Where("book_id = ? AND deleted_at IS NULL", bookid).Find(&booking)
	if (models.Booking{}) == booking {
		log.Println("Booking is not available")
		json.NewEncoder(w).Encode("Booking is not available")
	}
	db.Delete(&booking)
}

func Login(w http.ResponseWriter, r *http.Request) {
	db := CreateDbConnection()
	defer CloseConnection(db)
	var u, v models.User
	json.NewDecoder(r.Body).Decode(&u)
	db.Where("username=?", u.Username).Find(&v)
	var resp map[string]interface{}
	if (models.User{}) == v {
		log.Println("User Doesn't Exists")
		resp = map[string]interface{}{"status": false, "message": "User Doesn't Exists"}
	}
	resp = getToken(u, v)
	json.NewEncoder(w).Encode(resp)
}
func getToken(u, v models.User) map[string]interface{} {

	err := bcrypt.CompareHashAndPassword([]byte(v.Password), []byte(u.Password))
	if err != nil {
		log.Println("Not authenticated")
		resp := map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp
	}
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()
	fmt.Println(v)
	tk := &Claims{
		Username: v.Username,
		Name:     v.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		Role: v.Role,
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		log.Panicln("Error in creating Token")
		resp := map[string]interface{}{"status": false, "message": "Error in creating Token"}
		return resp
	}
	var resp = map[string]interface{}{"status": true, "message": "logged in"}
	resp["token"] = tokenString //Store the token in the response
	resp["user"] = v
	return resp
}

func JwtRegularVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var header = r.Header.Get("x-access-token")

		header = strings.TrimSpace(header)

		if header == "" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Missing auth token")
			return
		}
		tk := &models.Token{}

		_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Message:" + err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func JwtAdminVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var header = r.Header.Get("x-access-token")

		header = strings.TrimSpace(header)

		if header == "" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Missing auth token")
			return
		}
		tk := &models.Token{}

		_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Message:" + err.Error())
			return
		}
		fmt.Println(*&tk.Role)
		if *&tk.Role != "admin" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Message: Forbidden A")
			return
		}
		ctx := context.WithValue(r.Context(), "user", tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
