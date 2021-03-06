package app

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

var (
	// ErrBadContentType is returned when the server gets an unexpected
	// content type.
	ErrBadContentType = errors.New("wrong content type")
)

// UploadSavePhoto saves an uploaded file to the filesystem.
func UploadSavePhoto(f io.Reader, name string, mediaConfig MediaConfig) (string, string, error) {
	newUUID, _ := uuid.NewV4()
	newID := newUUID.String()

	ext := filepath.Ext(name)
	path := mediaConfig.Path + newID + ext
	thumbPath := mediaConfig.ThumbnailsPath + "thumb_" + newID + ext

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return "", "", err
	}

	contentType := detectContentType(fileBytes)
	if !okContentType(contentType) {
		return "", "", ErrBadContentType
	}

	err = writeFile(path, fileBytes, 0644)
	if err != nil {
		return "", "", err
	}
	err = genThumbnail(thumbPath, fileBytes, contentType)
	if err != nil {
		return "", "", err
	}
	return newID, path, err
}

func okContentType(contentType string) bool {
	return contentType == "image/png" || contentType == "image/jpeg" || contentType == "image/gif"
}

// detectContentType from
// https://golangcode.com/get-the-content-type-of-file/
func detectContentType(fb []byte) string {
	// Only the first 512 bytes are used to sniff the content type.
	// Use the net/http package's handy DetectContentType function.
	// Always seems to return a valid content-type by returning
	// "application/octet-stream" if no others seemed to match.
	return http.DetectContentType(fb[:512])
}

func writeFile(p string, f []byte, fmode os.FileMode) error {
	return os.WriteFile(p, f, fmode)
}

// RemoveFile removes a file from the filesystem.
func RemoveFile(f string) error {
	log.Print("Removing item from filesystem: ", f)
	return os.Remove(f)
}
