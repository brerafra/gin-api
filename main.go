package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// album represents data about a record album
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

//albums slice to seed records album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

//User struct
type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email"`
}

//postAlbums adds an album from josn received in the request body
func postAlbums(c *gin.Context) {
	var newAlbum album
	//Call Bindjson to bind teh received JSON to newAlbum,
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	fmt.Println(newAlbum)

	//Add the new album to the slice
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

//getAlbumByID locates the album whose ID value matches the id
//parameter sent by the client, then returns that album as a response
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	//loop over the list of albums, looking for an album
	//whose id value matches the parameter
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
}

//getAlbums responds with the list  of all albums as Json
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func main() {
	fmt.Println("gorm-gin usos")
	handleRequsts()
}

func handleRequsts() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.GET("/albums/:id", getAlbumByID)

	router.GET("/users", allUsers)
	router.POST("/users", newUser)
	initialMigration()

	router.Run("localhost:8080")
}

func allUsers(c *gin.Context) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var users []User

	db.Find(&users)
	fmt.Println("{}", users)
	c.IndentedJSON(http.StatusOK, users)
}

func newUser(c *gin.Context) {
	var newUser User
	//Call Bindjson to bind teh received JSON to newAlbum,
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	fmt.Println(newUser)
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.Create(&newUser)

	c.IndentedJSON(http.StatusCreated, newUser)
}

func initialMigration() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("falied to conect database")
	}
	defer db.Close()

	//Migrate the schema
	db.AutoMigrate(&User{})
}
