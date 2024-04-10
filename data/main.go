package data

import (
	"backend/exceptions"
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

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

func CloseDB() {
	err := db.Close()
	if err != nil {
		log.Fatalln("Error closing the database connection")
	}
}

func startTransaction() (*sql.Tx, *exceptions.DataPackageError) {
	// On ouvre une transaction BDD
	tx, errTx := db.Begin()
	// Si erreur, on plante
	if errTx != nil {
		// Gestion d'erreur
		return nil, &exceptions.DataPackageError{Message: "Unable to start transaction", Code: exceptions.SQL_ERROR_TRANS_BEGIN}
	}

	return tx, nil
}

func closeTransaction(tx *sql.Tx) *exceptions.DataPackageError {
	// Commit la transaction
	errTx := tx.Commit()
	// Gestion erreur
	if errTx != nil {
		// Autre type d'erreur
		return &exceptions.DataPackageError{Message: "Unable to commit transaction", Code: exceptions.SQL_ERROR_TRANS_STOP}
	}
	return nil
}

func manageSqlError(errEx error, tx *sql.Tx) *exceptions.DataPackageError {
	// Gestion erreur
	if errEx != nil {
		// Vérifier si c'est une erreur MySQL
		var mysqlErr *mysql.MySQLError
		if errors.As(errEx, &mysqlErr) {
			// Vérifier si c'est une erreur de clé dupliquée
			if mysqlErr.Number == 1062 {
				// Retour de l'erreur de duplication de clé
				return &exceptions.DataPackageError{Message: "Duplicate key insertion", Code: exceptions.SQL_ERROR_DUPLICATE}
			} else {
				// Autre type d'erreur MySQL
				return &exceptions.DataPackageError{Message: "SQL error", Code: exceptions.SQL_ERROR_LAMBDA}
			}
		} else {
			// Autre type d'erreur
			// Tentative de rollback
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return &exceptions.DataPackageError{Message: "Unable to rollback", Code: exceptions.SQL_ERROR_TRANS_ROLLBACK}
			}
			return &exceptions.DataPackageError{Message: "Internal error", Code: exceptions.ERROR_LAMBDA}
		}
	}
	return nil
}
