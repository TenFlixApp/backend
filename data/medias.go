package data

import (
	"backend/exceptions"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type MediaPreview struct {
	Titre       string   `json:"titre"`
	Description string   `json:"description"`
	UUIDMedia   string   `json:"uuid"`
	Createur    UserInfo `json:"createur"`
}

func parseMedias(rows *sql.Rows, createur *UserInfo) ([]MediaPreview, error) {
	medias := make([]MediaPreview, 0)
	for rows.Next() {
		media := MediaPreview{}
		var err error
		if createur == nil {
			crea := UserInfo{}
			err = rows.Scan(&media.Titre, &media.Description, &media.UUIDMedia, &crea.Nom, &crea.Prenom, &crea.UUIDAvatar)
			media.Createur = crea
		} else {
			err = rows.Scan(&media.Titre, &media.Description, &media.UUIDMedia)
			media.Createur = *createur
		}
		if err != nil {
			return nil, errors.New("errors when getting media's info")
		}
		medias = append(medias, media)
	}
	return medias, nil
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

	return nil
}

func GetMedias(ids []string) ([]MediaPreview, error) {
	param := fmt.Sprintf("'%s'", strings.Join(ids, "','"))
	query := fmt.Sprintf("SELECT m.titre, m.description, m.uuid_media, u.nom, u.prenom, u.uuid_avatar FROM medias m JOIN users u ON (m.id_createur = u.id_user) WHERE m.uuid_media IN (%s)", param)
	rows, err := db.Query(query)
	if err != nil {
		return nil, errors.New("unable to get medias")
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	return parseMedias(rows, nil)
}

func GetMediaFromCreator(idCreator int64) ([]MediaPreview, error) {
	createur, err := GetUserInfo(idCreator)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT titre, description, uuid_media FROM medias WHERE id_createur = ? ORDER BY updated_at DESC", idCreator)
	if err != nil {
		return nil, errors.New("unable to get creator's medias")
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	return parseMedias(rows, createur)
}

func SearchMedia(searchTerm string) ([]MediaPreview, error) {
	rows, err := db.Query(`SELECT m.titre, m.description, m.uuid_media, u.nom, u.prenom, u.uuid_avatar FROM medias m JOIN users u ON (m.id_createur = u.id_user) WHERE titre LIKE ?`, "%"+searchTerm+"%")
	if err != nil {
		return nil, errors.New("unable to search media")
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	return parseMedias(rows, nil)
}
