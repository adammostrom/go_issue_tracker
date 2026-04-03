package main

import (
	"flag"
	"fmt"
	"issuetracker/internal/cli"
	database "issuetracker/internal/database"
	router "issuetracker/internal/router"
	services "issuetracker/internal/services"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/lib/pq" // Import pq for PostgreSQL driver
)

// Main does only one job: start the server and connect routes to handlers.
// Only does initialization: database connection, routing, and starting the HTTP server.
func main() {

	// Initate and load the database schema
	db := database.InitDB()

	// Initiate the connection to the database (read/write) -> layer between service and db
	db_connection := database.NewDatabaseConnection(db)

	// Create a service struct that delegates to the database for read/write
	issueService := services.NewIssueService(db_connection)

	switch os.Args[1] {
	case "start":
		startCmd(os.Args[2:], issueService)
	case "issue":
		issueCmd(os.Args[2:], issueService)
	}

}

func issueCmd(args []string, issueService *services.IssueService) {

	cli := cli.NewCLI(issueService)

	if len(args) < 1 {
		fmt.Println("expected subcommand")
		return
	}
	switch args[0] {
	case "list":
		cli.GetIssues()
	case "show":
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("invalid id: %s\n", args[0])
			return
		}
		issue, err := cli.GetIssue(id)
		if err != nil {
			return
		}
		fmt.Printf("issue found: %s\n", issue.Title)
	}

}

// Start the HTTP server
func startCmd(args []string, issueService *services.IssueService) {
	fs := flag.NewFlagSet("start", flag.ExitOnError)

	// Create a router to delegate requests to the server
	r := router.NewRouter(issueService)

	fs.Parse(args)

	// set up the HTTP server
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080", // Add configuration if want to change ports etc.
		Handler: mux,
	}

	// Delegates all HTTP requests to /issues* to the IssuesHandler
	mux.HandleFunc("/issues", r.AllRouting)  // for /issues exact (list, create)
	mux.HandleFunc("/issues/", r.AllRouting) // for /issues/{id} (single issue)

	log.Println("Connected successfully — server starting...")

	// Serve the HTML frontend
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	// Keep server running on port 8080
	log.Println("Server running at http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
