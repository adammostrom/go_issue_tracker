package main

import (
	// Provides SQL database functionality
	// For printing messages
	// For printing messages

	//"issuetracker/internal/db"
	handlers "issuetracker/internal/backend_logic"
	db "issuetracker/internal/database_layer"
	"issuetracker/internal/models"
	"issuetracker/internal/services"
	"log"
	"net/http"

	_ "github.com/lib/pq" // Import pq for PostgreSQL driver
)

// Main does only one job: start the server and connect routes to handlers.
// Only does initialization: database connection, routing, and starting the HTTP server.
// Doesnt do actual “business work”.
func main() {

	// 2026-03-12 -> CREATE NEW ISSUE:

	s := handlers.IssueEndpoint{}

	// TEMPORARY STORAGE
	var issues []models.Issue
	var nextID = 1

	var issue = s.CreateNewIssue("First Test", "just testing", int64(nextID))

	issues = append(issues, issue)

	service := services.IssueService{
		Issues: &issues,
	}

	http.HandleFunc("/issues", service.MainRouter)

	log.Println("Connected successfully — server starting...")

	// Serve the HTML frontend
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Keep server running on port 8080
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	// DATABASE

	/*
		host=/var/run/postgresql forces a Unix socket.
		peer auth sees you’re adam on the system and lets you in without a password.
		No password= needed.
		This is the easiest and safe for local dev.
	*/
	connStr := "user=adam dbname=servicetestess host=/var/run/postgresql sslmode=disable"

	// connect using the datapase package, which contains the function OpenDB in db_connect.go
	db_connect, err := db.OpenDB(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db_connect.Close()

	// Add functionality for frontend later.
	//getAllDevices(db_connect)
	//getAllDeviceTypes((db_connect))
}

/* func getAllDevices(db_connect *sql.DB) {
	deviceStore := db.NewDeviceStore(db_connect)

	// Important to not ignore the error, if GetAll fails, devices will default to "nil" (none)
	// And this we need/want to handle.
	devices, err := deviceStore.GetAll()
	if err != nil {
		log.Println("GetAll failed:", err)
		return
	}

	for _, d := range devices {
		fmt.Printf("DeviceID: %d DeviceType: %s \n", d.DeviceID, d.DeviceType)
	}

}

func getAllDeviceTypes(db_connect *sql.DB) {
	deviceTypeStore := db.NewDeviceTypeStore(db_connect)

	// Important to not ignore the error, if GetAll fails, devices will default to "nil" (none)
	// And this we need/want to handle.
	devices, err := deviceTypeStore.GetAll()
	if err != nil {
		log.Println("GetAll failed:", err)
		return
	}

	for _, dt := range devices {
		fmt.Printf("DeviceType: %s \n", dt.DeviceType)
	}

}
*/
