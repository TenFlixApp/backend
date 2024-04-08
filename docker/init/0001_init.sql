CREATE TABLE IF NOT EXISTS users (
    id_user INT NOT NULL AUTO_INCREMENT,
    email VARCHAR(255) NOT NULL UNIQUE,
    nom VARCHAR(255) NOT NULL,
    prenom VARCHAR(255) NOT NULL,
    uuid_avatar TEXT DEFAULT '' COMMENT 'uuid avatar de l''utilisateur',
    PRIMARY KEY (id_user)
);

CREATE TABLE IF NOT EXISTS medias (
    id_media INT NOT NULL AUTO_INCREMENT,
    id_createur INT NOT NULL,
    titre VARCHAR(255) NOT NULL,
    uuid_media TEXT DEFAULT '' COMMENT 'uuid du media',
    description TEXT DEFAULT '',
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (id_media),
    CONSTRAINT fk_medias_createur
        FOREIGN KEY (id_createur) REFERENCES users (id_user)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS playlists (
    id_playlist INT NOT NULL AUTO_INCREMENT,
    id_createur INT NOT NULL,
    titre VARCHAR(255) DEFAULT 'Playlist',
    PRIMARY KEY (id_playlist),
    CONSTRAINT fk_playlists_createur
        FOREIGN KEY (id_createur) REFERENCES users (id_user)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS playlists_medias (
    id_playlist INT NOT NULL,
    id_media INT NOT NULL,
    CONSTRAINT fk_playlists_medias_playlist
        FOREIGN KEY (id_playlist) REFERENCES playlists (id_playlist)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
    CONSTRAINT fk_playlists_medias_media
        FOREIGN KEY (id_media) REFERENCES medias (id_media)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS commentaires (
    id_commentaire INT NOT NULL AUTO_INCREMENT,
    id_createur INT NOT NULL,
    id_media INT NOT NULL,
    id_parent INT DEFAULT NULL,
    commentaire TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (id_commentaire),
    CONSTRAINT fk_commentaires_createur
        FOREIGN KEY (id_createur) REFERENCES users (id_user)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
    CONSTRAINT fk_commentaires_media
        FOREIGN KEY (id_media) REFERENCES medias (id_media)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
    CONSTRAINT fk_commentaires_parent
        FOREIGN KEY (id_parent) REFERENCES commentaires (id_commentaire)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
);