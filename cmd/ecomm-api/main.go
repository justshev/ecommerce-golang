package main

import (
	"log"

	"github.com/mshevaatallah/ecomm/db"
)

func main() {
db, err := db.NewDatabase()
if err != nil {
	log.Fatalf("Error %v", err)

}
defer db.Close()
log.Println("Connected to database")
}