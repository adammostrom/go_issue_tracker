package main

import (
	"database/sql"
	"flag"
	"fmt"
	"issuetracker/internal/cli"
	database "issuetracker/internal/database"
	router "issuetracker/internal/router"
	services "issuetracker/internal/services"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/lib/pq" // Import pq for PostgreSQL driver
)

// Main does only one job: start the server and connect routes to handlers.
// Only does initialization: database connection, routing, and starting the HTTP server.
func main() {

	if len(os.Args) < 2 {
		fmt.Println("Exptected commands: start, issue, init")
		return
	}

	cmd := os.Args[1]

	var db *sql.DB
	var err error

	switch cmd {
	case "init":
		_, err := database.InitDB()
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		// all other commands require existing db
		db, err = database.OpenDB()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	// Initate and load the database schema
	// Initiate the connection to the database (read/write) -> layer between service and db
	db_connection := database.NewDatabaseConnection(db)

	// Create a service struct that delegates to the database for read/write
	issueService := services.NewIssueService(db_connection)

	cli := cli.NewCLI(issueService)

	cmds := cli.BuildCommands()

	switch os.Args[1] {
	case "init":
		fmt.Println("Database already initiated")
	case "start":
		startCmd(os.Args[2:], issueService)
	case "issue":
		cli.Run(cmds, os.Args[2:])
	default:
		fmt.Println("No valid command provided") // TODO: Show subcommands (make function)

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
