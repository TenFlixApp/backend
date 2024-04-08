package exceptions

// Erreurs SQL Gérées
const (
	SQL_ERROR_TRANS_COMMIT   = 5 // Erreur lors du rollback de la transaction
	SQL_ERROR_TRANS_ROLLBACK = 4 // Erreur lors du rollback de la transaction
	SQL_ERROR_TRANS_BEGIN    = 3 // Erreur lors du lancement de la transaction SQL
	SQL_ERROR_DUPLICATE      = 2 // Duplication de clé primaire
	SQL_ERROR_LAMBDA         = 1 // Erreur SQL par défaut
	ERROR_LAMBDA             = 0 // Erreur par défaut
)

// Type d'erreur remontable par le package
type DataPackageError struct {
	Message string
	Code    int
}
