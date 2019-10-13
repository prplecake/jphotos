package db

import (
	"fmt"
	"time"
)

// A Photo is a view into a photo
type Photo struct {
	ID, Caption, Location string
	Added                 time.Time
}

// AddPhoto adds a photo to the database
func (pg *PGStore) AddPhoto(p Photo, albumID string) error {
	_, err := pg.Query(
		"INSERT INTO photos (id, caption, location) "+
			"VALUES ($1, $2, $3)",
		p.ID, p.Caption, p.Location)
	if err != nil {
		return fmt.Errorf("AddPhoto: %w", err)
	}

	_, err = pg.Query(
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
	_, err := pg.Query("DELETE FROM photos WHERE id = $1",
		id)
	if err != nil {
		return fmt.Errorf("DeletePhotoByID: %w", err)
	}

	return nil
}
