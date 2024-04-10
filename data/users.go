package data

import (
	"backend/exceptions"
	"database/sql"
	"errors"
)

type UserInfo struct {
	Nom        string `json:"nom"`
	Prenom     string `json:"prenom"`
	UUIDAvatar string `json:"avatar"`
}

func RegisterUser(email string, nom string, prenom string) (err *exceptions.DataPackageError) {
	tx, errData := startTransaction()
	if errData != nil {
		return errData
	}

	// Insertion de l'utilisateur
	_, errEx := tx.Exec(`INSERT INTO users (email, nom, prenom) VALUES (?, ?, ?)`, email, nom, prenom)

	errData = manageSqlError(errEx, tx)
	if errData != nil {
		return errData
	}

	errData = closeTransaction(tx)
	if errData != nil {
		return errData
	}

	// Retourne pas d'erreur
	return nil
}

func GetUserIDByEmail(email string) (int, error) {
	var userID int

	// Exécuter la requête SQL pour récupérer l'ID de l'utilisateur avec l'email spécifié
	err := db.QueryRow("SELECT id_user FROM users WHERE email = ?", email).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("user not found")
		}
		return 0, err
	}

	return userID, nil
}

func ChangeAvatar(uuid string, idUser int) *exceptions.DataPackageError {
	tx, errData := startTransaction()
	if errData != nil {
		return errData
	}

	// Insertion de l'utilisateur
	_, errEx := tx.Exec(`UPDATE users SET uuid_avatar = ? WHERE email = ?;`, uuid, idUser)

	errData = manageSqlError(errEx, tx)
	if errData != nil {
		return errData
	}

	errData = closeTransaction(tx)
	if errData != nil {
		return errData
	}

	// Retourne pas d'erreur
	return nil
}

func GetUserInfo(id int64) (*UserInfo, error) {
	user := new(UserInfo)

	err := db.QueryRow("SELECT nom, prenom, uuid_avatar FROM users WHERE id_user = ?", id).Scan(&user.Nom, &user.Prenom, &user.UUIDAvatar)
	if err != nil {
		return nil, errors.New("unable to get user info")
	}

	return user, nil
}
