package db

import (
	"fmt"
	"log"
	"time"
)

// A Photo is a view into a photo
type Photo struct {
	UUID, Caption, Location string
	Added                   time.Time
	ID                      int
}

// AddPhoto adds a photo to the database
func (pg *PGStore) AddPhoto(p Photo, albumUUID string) error {
	err := pg.Exec(
		"INSERT INTO photos (uuid, caption, location) "+
			"VALUES ($1, $2, $3)",
		p.UUID, p.Caption, p.Location)
	if err != nil {
		return fmt.Errorf("AddPhoto: %w", err)
	}

	err = pg.Exec(
		"INSERT INTO album_photos (photo, album) "+
			"VALUES ($1, $2)",
		p.UUID, albumUUID)
	if err != nil {
		return fmt.Errorf("AddPhoto: %w", err)
	}
	return nil
}

// GetPhotoByUUID returns a photo object for that ID.
func (pg *PGStore) GetPhotoByUUID(uuid string) (*Photo, error) {
	rows, err := pg.Query(
		"SELECT caption, location, added, id "+
			"FROM photos WHERE uuid = $1",
		uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, ErrNotFound
	}

	var (
		caption, location string
		added             time.Time
		id                int
	)

	err = rows.Scan(&caption, &location, &added, &id)
	if err != nil {
		return nil, err
	}

	return &Photo{
		UUID:     uuid,
		Caption:  caption,
		Location: location,
		Added:    added,
		ID:       id,
	}, nil
}

// DeletePhotoByUUID deletes the photo at the provided ID
func (pg *PGStore) DeletePhotoByUUID(uuid string) error {
	txn, err := pg.conn.Begin()
	if err != nil {
		log.Printf("Currently there are %d connections.", pg.conn.Stats().OpenConnections)
		return fmt.Errorf("DeletePhotoByID/Begin(): %w", err)
	}
	_, err = txn.Exec("DELETE FROM photos WHERE uuid = $1", uuid)
	if err != nil {
		return fmt.Errorf("DeletePhotoByID/Exec(): %w", err)
	}
	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("DeletePhotoByID/Commit(): %w", err)
	}

	return nil
}

// UpdatePhotoCaptionByUUID updates the photo's caption
func (pg *PGStore) UpdatePhotoCaptionByUUID(uuid, newCaption string) error {

	err := pg.Exec(
		"UPDATE photos SET caption = $1 WHERE uuid = $2",
		newCaption, uuid)
	if err != nil {
		return fmt.Errorf("UpdatePhotoCaption: %w", err)
	}
	log.Print("Photo caption updated.")
	return nil
}

// UpdatePhotoAlbum changes the album a photo belongs to
func (pg *PGStore) UpdatePhotoAlbum(photoUUID, albumUUID string) error {
	return pg.Exec(
		"UPDATE album_photos SET album = $2 WHERE photo = $1",
		photoUUID, albumUUID)
}

// GetAlbumUUIDByPhotoUUID returns the album slug a photo belongs to
func (pg *PGStore) GetAlbumUUIDByPhotoUUID(photoUUID string) (string, error) {
	rows, err := pg.Query(
		"SELECT a.uuid FROM albums AS a "+
			"INNER JOIN album_photos AS ap ON ap.album = a.uuid "+
			"WHERE ap.photo = $1",
		photoUUID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", ErrNotFound
	}

	var albumUUID string
	err = rows.Scan(&albumUUID)
	if err != nil {
		return "", err
	}

	return albumUUID, nil
}
