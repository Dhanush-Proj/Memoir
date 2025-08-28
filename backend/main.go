package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Blog struct represents a blog entry.
type Blog struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Image   string `json:"image"`
}

// Global variable for the database connection.
var db *sql.DB

// getBlogs handles GET requests to /blogs.
func getBlogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var blogs []Blog
	rows, err := db.Query("SELECT id, date, title, content, image FROM blogs")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var blog Blog
		err := rows.Scan(&blog.ID, &blog.Date, &blog.Title, &blog.Content, &blog.Image)
		if err != nil {
			log.Fatal(err)
		}
		blogs = append(blogs, blog)
	}
	json.NewEncoder(w).Encode(blogs)
}

// getBlog handles GET requests to /blogs/{id}.
func getBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var blog Blog
	row := db.QueryRow("SELECT id, date, title, content, image FROM blogs WHERE id = ?", params["id"])
	err := row.Scan(&blog.ID, &blog.Date, &blog.Title, &blog.Content, &blog.Image)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Blog not found", http.StatusNotFound)
			return
		}
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(blog)
}

// createBlog handles POST requests to /blogs.
func createBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var blog Blog
	_ = json.NewDecoder(r.Body).Decode(&blog)
	result, err := db.Exec("INSERT INTO blogs(date, title, content, image) VALUES(?, ?, ?, ?)", blog.Date, blog.Title, blog.Content, blog.Image)
	if err != nil {
		log.Fatal(err)
	}
	lastID, _ := result.LastInsertId()
	blog.ID = string(lastID)
	json.NewEncoder(w).Encode(blog)
}

// EditBlog handles PUT requests to /blogs/{id}.
func EditBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var blog Blog
	_ = json.NewDecoder(r.Body).Decode(&blog)
	_, err := db.Exec("UPDATE blogs SET date = ?, title = ?, content = ?, image = ? WHERE id = ?", blog.Date, blog.Title, blog.Content, blog.Image, params["id"])
	if err != nil {
		log.Fatal(err)
	}
	blog.ID = params["id"]
	json.NewEncoder(w).Encode(blog)
}

// DeleteBlog handles DELETE requests to /blogs/{id}.
func DeleteBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	_, err := db.Exec("DELETE FROM blogs WHERE id = ?", params["id"])
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Blog deleted successfully"})
}

func main() {
	// Open a connection to the SQLite database.
	var err error
	db, err = sql.Open("sqlite3", "./blogs.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the "blogs" table if it doesn't exist.
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS blogs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date TEXT,
        title TEXT,
        content TEXT,
        image TEXT
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	router := mux.NewRouter()
	// In your main() function:
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"}) // Allows requests from any origin
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	// Define the API routes.
	router.HandleFunc("/blogs", getBlogs).Methods("GET")
	router.HandleFunc("/blogs/{id}", getBlog).Methods("GET")
	router.HandleFunc("/blogs", createBlog).Methods("POST")
	router.HandleFunc("/blogs/{id}", EditBlog).Methods("PUT")
	router.HandleFunc("/blogs/{id}", DeleteBlog).Methods("DELETE")

	// Start the server.
	// Old code:
	// log.Fatal(http.ListenAndServe(":8080", router))

	// New code with CORS:
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
