package data

import (
	"backend/exceptions"

	"github.com/go-sql-driver/mysql"
)

func RegisterUser(email string, nom string, prenom string) (err *exceptions.DataPackageError) {
	// On ouvre une transaction BDD
	tx, errTx := db.Begin()
	// Si erreur, on plante
	if errTx != nil {
		// Gestion d'erreur
		return &exceptions.DataPackageError{Message: "Unable to start transaction", Code: exceptions.SQL_ERROR_LAMBDA}
	}

	// Insertion de l'utilisateur
	_, errEx := tx.Exec(`INSERT INTO users (email, nom, prenom) VALUES (?, ?, ?)`, email, nom, prenom)
	// Gestion erreur
	if errEx != nil {
		// Vérifier si c'est une erreur MySQL
		if mysqlErr, ok := errEx.(*mysql.MySQLError); ok {
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

	// Commit la transaction
	errTx = tx.Commit()
	// Gestion erreur
	if errTx != nil {
		// Autre type d'erreur
		return &exceptions.DataPackageError{Message: "Unable to commit transaction", Code: exceptions.ERROR_LAMBDA}
	}

	// Retourne pas d'erreur
	return nil
}
