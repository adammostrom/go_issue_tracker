package main

import (
	// Provides SQL database functionality
	// For printing messages
	// For printing messages

	"database/sql"
	"fmt"
	"log"
	"net/http"
	"servicetestess/database"

	_ "github.com/lib/pq" // Import pq for PostgreSQL driver
)

func main() {

	/*
		host=/var/run/postgresql forces a Unix socket.
		peer auth sees you’re adam on the system and lets you in without a password.
		No password= needed.
		This is the easiest and safe for local dev.
	*/
	connStr := "user=adam dbname=servicetestess host=/var/run/postgresql sslmode=disable"

	//ADASDASDAS
	// connect using the datapase package, which contains the function OpenDB in db_connect.go
	db, err := database.OpenDB(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Connected successfully — server starting...")

	// Add functionality for frontend later.
	getAllDevices(db)
	getAllDeviceTypes((db))

	// Serve the HTML frontend
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Keep server running on port 8080
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getAllDevices(db *sql.DB) {
	deviceStore := database.NewDeviceStore(db)

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

func getAllDeviceTypes(db *sql.DB) {
	deviceTypeStore := database.NewDeviceTypeStore(db)

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
