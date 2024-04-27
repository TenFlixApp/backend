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
	tx, errTx := db.Begin()
	if errTx != nil {
		return nil, &exceptions.DataPackageError{Message: "Unable to start transaction", Code: exceptions.SQL_ERROR_TRANS_BEGIN}
	}

	return tx, nil
}

func closeTransaction(tx *sql.Tx) *exceptions.DataPackageError {
	errTx := tx.Commit()
	if errTx != nil {
		return &exceptions.DataPackageError{Message: "Unable to commit transaction", Code: exceptions.SQL_ERROR_TRANS_STOP}
	}
	return nil
}

func manageSqlError(errEx error, tx *sql.Tx) *exceptions.DataPackageError {
	if errEx != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(errEx, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return &exceptions.DataPackageError{Message: "Duplicate key insertion", Code: exceptions.SQL_ERROR_DUPLICATE}
			} else {
				return &exceptions.DataPackageError{Message: "SQL error", Code: exceptions.SQL_ERROR_LAMBDA}
			}
		} else {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return &exceptions.DataPackageError{Message: "Unable to rollback", Code: exceptions.SQL_ERROR_TRANS_ROLLBACK}
			}
			return &exceptions.DataPackageError{Message: "Internal error", Code: exceptions.ERROR_LAMBDA}
		}
	}
	return nil
}
