package db

import (
	"fmt"
	"log"
	"time"
)

// A Photo is a view into a photo
type Photo struct {
	ID, Caption, Location string
	Added                 time.Time
}

// AddPhoto adds a photo to the database
func (pg *PGStore) AddPhoto(p Photo, albumID string) error {
	err := pg.Exec(
		"INSERT INTO photos (id, caption, location) "+
			"VALUES ($1, $2, $3)",
		p.ID, p.Caption, p.Location)
	if err != nil {
		return fmt.Errorf("AddPhoto: %w", err)
	}

	err = pg.Exec(
		"INSERT INTO album_photos (photo, album) "+
			"VALUES ($1, $2)",
		p.ID, albumID)
	if err != nil {
		return fmt.Errorf("AddPhoto: %w", err)
	}
	return nil
}

// GetPhotoByID returns a photo object for that ID.
func (pg *PGStore) GetPhotoByID(id string) (*Photo, error) {
	rows, err := pg.Query(
		"SELECT caption, location, added "+
			"FROM photos WHERE id = $1",
		id)
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
	)

	err = rows.Scan(&caption, &location, &added)
	if err != nil {
		return nil, err
	}

	return &Photo{
		ID:       id,
		Caption:  caption,
		Location: location,
		Added:    added,
	}, nil
}

// DeletePhotoByID deletes the photo at the provided ID
func (pg *PGStore) DeletePhotoByID(id string) error {
	txn, err := pg.conn.Begin()
	if err != nil {
		log.Printf("Currently there are %d connections.", pg.conn.Stats().OpenConnections)
		return fmt.Errorf("DeletePhotoByID/Begin(): %w", err)
	}
	_, err = txn.Exec("DELETE FROM photos WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("DeletePhotoByID/Exec(): %w", err)
	}
	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("DeletePhotoByID/Commit(): %w", err)
	}

	return nil
}

// UpdatePhotoCaption updates the photo's caption
func (pg *PGStore) UpdatePhotoCaption(id, newCaption string) error {
	err := pg.Exec(
		"UPDATE photos SET caption = $1 WHERE id = $2",
		newCaption, id)
	if err != nil {
		return fmt.Errorf("UpdatePhotoCaption: %w", err)
	}
	log.Print("Photo caption updated.")
	return nil
}

// GetPhotoAlbum returns the album slug a photo belongs to
func (pg *PGStore) GetPhotoAlbum(photoID string) (string, error) {
	rows, err := pg.Query(
		"SELECT a.slug FROM albums AS a "+
			"INNER JOIN album_photos AS ap ON ap.album = a.id "+
			"WHERE ap.photo = $1",
		photoID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", ErrNotFound
	}

	var albumID string
	err = rows.Scan(&albumID)
	if err != nil {
		return "", err
	}

	return albumID, nil
}
