package data

import (
	"backend/exceptions"
	"database/sql"
	"errors"
	"fmt"
)

type Playlist struct {
	ID    int64  `json:"id"`
	Titre string `json:"titre"`
}

func CreatePlaylist(idCreateur int, titre string) *exceptions.DataPackageError {
	tx, errData := startTransaction()
	if errData != nil {
		return errData
	}

	_, errEx := tx.Exec(`INSERT INTO playlists (titre, id_createur) VALUES (?, ?)`, titre, idCreateur)

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

func DeletePlaylist(id int64) *exceptions.DataPackageError {
	tx, errData := startTransaction()
	if errData != nil {
		return errData
	}

	_, errEx := tx.Exec(`DELETE FROM playlists WHERE id_playlist = ?`, id)

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

func AddMediaToPlaylist(idPlaylist int, uuidMedia string) *exceptions.DataPackageError {
	tx, errData := startTransaction()
	if errData != nil {
		return errData
	}

	fmt.Print(idPlaylist)

	_, errEx := tx.Exec(`INSERT INTO playlists_medias (id_playlist, id_media) VALUES (?, (SELECT id_media FROM medias WHERE uuid_media = ?))`, idPlaylist, uuidMedia)

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

func DeleteMediaFromPlaylist(idPlaylist int, uuidMedia string) *exceptions.DataPackageError {
	tx, errData := startTransaction()
	if errData != nil {
		return errData
	}

	_, errEx := tx.Exec(`DELETE FROM playlists_medias WHERE id_playlist = ? AND id_media = (SELECT id_media FROM medias WHERE uuid_media = ?)`, idPlaylist, uuidMedia)

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

func GetPlaylistsFromUser(idUser int) ([]Playlist, *exceptions.DataPackageError) {

	rows, err := db.Query(`SELECT id_playlist, titre FROM playlists WHERE id_createur = ?`, idUser)
	if err != nil {
		return nil, &exceptions.DataPackageError{Message: "error when getting playlists", Code: 500}
	}

	playlists := make([]Playlist, 0)
	for rows.Next() {
		playlist := Playlist{}
		err := rows.Scan(&playlist.ID, &playlist.Titre)
		if err != nil {
			return nil, &exceptions.DataPackageError{Message: "error when getting playlists", Code: 500}
		}
		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func GetMediaFromPlaylist(idPlaylist int64) ([]MediaPreview, error) {
	rows, err := db.Query(`
		SELECT m.titre, m.description, m.uuid_media,
		u.nom, u.prenom, u.uuid_avatar
		FROM medias m
		JOIN users u ON (m.id_createur = u.id_user)
		JOIN playlists_medias pm ON (m.id_media = pm.id_media)
		WHERE pm.id_playlist = ?
	`, idPlaylist)
	if err != nil {
		return nil, errors.New("error when getting medias from playlist")
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	return parseMedias(rows, nil)
}
