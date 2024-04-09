package data

import (
	"backend/exceptions"
	"database/sql"
	"errors"
)

type MediaPreview struct {
	Titre     string `json:"titre"`
	UUIDMedia string `json:"uuid"`
}

func DeleteMedia(id int64) *exceptions.DataPackageError {
	tx, errData := startTransaction()
	if errData != nil {
		return errData
	}

	_, errEx := tx.Exec(`DELETE FROM medias WHERE id_media = ?`, id)

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

func CreateMedia(idCreateur int, titre string, uuid string, description string) *exceptions.DataPackageError {
	tx, errData := startTransaction()
	if errData != nil {
		return errData
	}

	_, errEx := tx.Exec(`INSERT INTO medias (id_createur, titre, uuid_media, description) VALUES (?, ?, ?, ?)`, idCreateur, titre, uuid, description)

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

func GetMediaFromCreator(idCreator int64) ([]MediaPreview, error) {
	var medias []MediaPreview

	// Exécuter la requête SQL pour récupérer les vidéos de l'utilisateur triées par date de mise à jour décroissante
	rows, err := db.Query("SELECT titre, uuid_media FROM medias WHERE id_createur = ? ORDER BY updated_at DESC", idCreator)
	if err != nil {
		return nil, errors.New("unable to get creator's medias")
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	// Parcourir les lignes de résultats et les ajouter au tableau
	for rows.Next() {
		var media MediaPreview
		if err := rows.Scan(&media.Titre, &media.UUIDMedia); err != nil {
			return nil, errors.New("errors when getting media's info")
		}
		medias = append(medias, media)
	}

	// Vérifier les erreurs lors du parcours des résultats
	if err := rows.Err(); err != nil {
		return nil, errors.New("errors when getting media's info")
	}

	return medias, nil
}
