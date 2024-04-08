package data

import (
	"database/sql"
	"log"
	"os"
)

// Variable driver BDD
var db *sql.DB

/**
* Connexion à la base de données
 */
func ConnectToDB() {
	var err error
	db, err = sql.Open("mysql", os.Getenv("DB_CONN_STRING"))
	if err != nil {
		log.Fatal("Unable to create DB handle", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the DB", err)
	}
	log.Println("Connected to the database")
}

/**
* Fermeture de la connexion BDD
 */
func CloseDB() {
	err := db.Close()
	if err != nil {
		log.Fatalln("Error closing the database connection")
	}
}
