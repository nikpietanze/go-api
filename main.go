package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}

var books []Book

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	requestID := c.Param("id")

	fmt.Println(requestID)

	for _, item := range books {
		if item.ID == requestID {
			c.IndentedJSON(http.StatusOK, item)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, &Book{})
}

func createBook(c *gin.Context) {
	var newBook Book

	if err := c.BindJSON(&newBook); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	_ = json.NewDecoder(c.Request.Body).Decode(&newBook)
	newBook.ID = strconv.Itoa(rand.Intn(10000000))
	books = append(books, newBook)
	
	c.IndentedJSON(http.StatusCreated, newBook)
}

func updateBook(c *gin.Context) {
	requestID := c.Param("id")

	for i, item := range books {
		if item.ID == requestID {
			books = append(books[:i], books[i+1:]...)
			
			var book Book
			if err := c.BindJSON(&book); err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}

			book.ID = requestID
			books = append(books, book)
			c.IndentedJSON(http.StatusOK, book)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, books)
}

func deleteBook(c *gin.Context) {
	requestID := c.Param("id")

	for index, item := range books {
		if item.ID == requestID {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	c.IndentedJSON(http.StatusOK, books)
}

func main() {
	router := gin.Default()
    
	books = append(books, Book{ID: "1", Isbn: "448743", Title: "Book One", Author: &Author{FirstName: "John", LastName: "Doe"}})
	books = append(books, Book{ID: "2", Isbn: "234908", Title: "Book Two", Author: &Author{FirstName: "Jane", LastName: "Doe"}})

	router.GET("/api/books", getBooks)
	router.GET("/api/books/:id", getBook)
	router.POST("/api/books", createBook)
	router.PUT("/api/books/:id", updateBook)
	router.DELETE("/api/books/:id", deleteBook)

	port := ":8000"
	log.Fatal(router.Run("localhost" + port))
}