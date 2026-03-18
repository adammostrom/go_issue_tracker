package main

import (
	database "issuetracker/internal/database"
	router "issuetracker/internal/router"
	services "issuetracker/internal/services"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/lib/pq" // Import pq for PostgreSQL driver
)

// Main does only one job: start the server and connect routes to handlers.
// Only does initialization: database connection, routing, and starting the HTTP server.
// Doesnt do actual “business work”.
func main() {

	// Initate and load the database schema
	db := database.InitDB()

	// Initiate the connection to the database (read/write) -> layer between service and db
	db_connection := database.NewDatabaseConnection(db)

	// Create a service struct that delegates to the database for read/write
	issueService := services.NewIssueService(db_connection)

	// Create a router to delegate requests to the server
	r := router.NewRouter(issueService)
	// set up the HTTP server
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080", // Add configuration if want to change ports etc.
		Handler: mux,
	}

	// Delegates all HTTP requests to /issues* to the IssuesHandler
	mux.HandleFunc("/issues", r.MainDelegator)  // for /issues exact (list, create)
	mux.HandleFunc("/issues/", r.MainDelegator) // for /issues/{id} (single issue)

	log.Println("Connected successfully — server starting...")

	// Serve the HTML frontend
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	// Keep server running on port 8080
	log.Println("Server running at http://localhost:8080")
	log.Fatal(server.ListenAndServe())

	// Add functionality for frontend later.
	//getAllDevices(db_connect)
	//getAllDeviceTypes((db_connect))
}
