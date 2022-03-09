package db

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/metal3d/go-slugify"
)

// An Album is a view into an album
type Album struct {
	Name, Slug, UUID string
}

var (

	// ErrAlbumExists is returned when the unique requirement of an
	// album name is violated
	ErrAlbumExists = errors.New("DB: Album exists")
	// ErrAlbumNameInvalid is returned when an album name is not valid.
	// Most likely to fire when a slug is blank. (jphotos#60)
	ErrAlbumNameInvalid = errors.New("DB: Album name invalid")
)

// AddAlbum adds an album to the albums table of the database
func (pg *PGStore) AddAlbum(name string) error {
	now := time.Now()
	slug := strings.ToLower(slugify.Marshal(name))
	log.Printf("Length of slug: %d", len(slug))
	if len(slug) == 0 {
		return ErrAlbumNameInvalid
	}
	err := pg.Exec("INSERT INTO albums (name, slug, created)"+
		"VALUES ($1, $2, $3)",
		name, slug, now)
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23505" {
			return ErrAlbumExists
		}
	}
	return err
}

// GetAlbums returns a list of all Albums
func (pg *PGStore) GetAllAlbums() ([]Album, error) {
	rows, err := pg.Query("SELECT uuid, name, slug FROM albums " +
		"ORDER BY name ASC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	albums := make([]Album, 0)

	for rows.Next() {
		var name, slug, id string
		err := rows.Scan(&id, &name, &slug)
		if err != nil {
			return nil, fmt.Errorf("GetAlbums: Couldn't scan: %w", err)
		}
		albums = append(albums, Album{name, slug, id})
	}
	return albums, nil
}

// GetAlbum returns an album, if it exists and matches the provided id
func (pg *PGStore) GetAlbumBySlug(slug string) (*Album, error) {
	rows, err := pg.Query("SELECT name, uuid FROM albums WHERE slug = $1", slug)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, ErrNotFound
	}

	var name, uuid string

	err = rows.Scan(&name, &uuid)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		panic("Database guarantee not met; multiples albums with same name")
	}

	return &Album{
		Name: name,
		Slug: slug,
		UUID: uuid,
	}, nil
}

// GetAlbumPhotos returns a list of all photos in an album
func (pg *PGStore) GetAlbumPhotosByUUID(uuid string) ([]Photo, error) {
	rows, err := pg.Query(
		"SELECT p.uuid, p.caption, p.location, p.added, p.id "+
			"FROM album_photos as ap "+
			"INNER JOIN photos as p ON p.uuid = ap.photo "+
			"INNER JOIN albums as a ON ap.album = a.uuid "+
			"WHERE ap.album = $1 ORDER BY p.id ASC",
		uuid)
	if err != nil {
		return nil, fmt.Errorf("GetAlbumPhotos: %w", err)
	}
	photos := make([]Photo, 0)

	for rows.Next() {
		var (
			uuid, caption, location string
			added                   time.Time
			id                      int
		)
		err = rows.Scan(&uuid, &caption, &location, &added, &id)
		if err != nil {
			return nil, fmt.Errorf("GetAlbumPhotos: %w", err)
		}
		photos = append(photos, Photo{uuid, caption, location, added, id})
	}
	return photos, nil
}

// GetPreviousAlbumPhoto returns the next photo in the album.
func (pg *PGStore) GetPreviousAlbumPhoto(albumID string, currentPhotoID int) string {
	rows, _ := pg.Query(
		"SELECT p.uuid "+
			"FROM album_photos AS ap "+
			"INNER JOIN photos AS p ON p.uuid = ap.photo "+
			"INNER JOIN albums AS a on ap.album = a.uuid "+
			"WHERE ap.album = $1 "+
			"AND p.id < $2 ORDER BY p.id DESC",
		albumID, currentPhotoID)
	defer rows.Close()

	if !rows.Next() {
		return ""
	}

	var (
		uuid string
	)

	_ = rows.Scan(&uuid)

	return uuid
}

// GetNextAlbumPhoto returns the next photo in the album.
func (pg *PGStore) GetNextAlbumPhoto(albumID string, currentPhotoID int) string {
	rows, _ := pg.Query(
		"SELECT p.uuid "+
			"FROM album_photos AS ap "+
			"INNER JOIN photos AS p ON p.uuid = ap.photo "+
			"INNER JOIN albums AS a on ap.album = a.uuid "+
			"WHERE ap.album = $1 "+
			"AND p.id > $2 ORDER BY p.id ASC",
		albumID, currentPhotoID)
	defer rows.Close()

	if !rows.Next() {
		return ""
	}

	var (
		uuid string
	)

	_ = rows.Scan(&uuid)

	return uuid
}

// GetFirstXPhotosFromAlbumByID returns the top X photos that belong to
// that album ID.
// The idea is that the first X photos will be used to make the album
// covers.
func (pg *PGStore) GetFirstXPhotosFromAlbumByID(albumID string, x int) ([]Photo, error) {
	rows, err := pg.Query(
		"SELECT p.uuid, p.location "+
			"FROM album_photos as ap "+
			"INNER JOIN photos as p ON p.uuid = ap.photo "+
			"INNER JOIN albums as a ON ap.album = a.uuid "+
			"WHERE ap.album = $1 ORDER BY p.id ASC "+
			"LIMIT $2",
		albumID, x)
	if err != nil {
		return nil, fmt.Errorf("GetFirstXPhotosFromAlbumByID: %w", err)
	}
	photos := make([]Photo, 0)

	for rows.Next() {
		var (
			uuid, caption, location string
			added                   time.Time
			id                      int
		)
		err = rows.Scan(&uuid, &location)
		if err != nil {
			return nil, fmt.Errorf("GetFirstXPhotosFromAlbumByID: %w", err)
		}
		photos = append(photos, Photo{uuid, caption, location, added, id})
	}
	return photos, nil
}

// GetAlbumSlugByUUID returns the slug matching the provided ID
func (pg *PGStore) GetAlbumSlugByUUID(uuid string) (string, error) {
	rows, err := pg.Query("SELECT slug FROM albums WHERE uuid = $1", uuid)
	if err != nil {
		return "", err
	}

	if !rows.Next() {
		return "", ErrNotFound
	}

	var slug string

	err = rows.Scan(&slug)
	if err != nil {
		return "", err
	}

	if rows.Next() {
		panic("Database guarantee not met; multiples albums with same name")
	}

	return slug, nil
}

// DeleteAlbumBySlug deletes the album, and all photos in it,
// matching the slug
func (pg *PGStore) DeleteAlbumBySlug(slug string) error {
	return pg.Exec("DELETE FROM albums WHERE slug = $1", slug)
}

// RenameAlbumByUUID renames an album
func (pg *PGStore) RenameAlbumByUUID(uuid, newName string) error {
	return pg.Exec(
		"UPDATE albums SET name = $1, slug = $2 "+
			"WHERE uuid = $3",
		newName, strings.ToLower(slugify.Marshal(newName)), uuid)
}
