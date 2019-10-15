package app

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
)

// UploadSavePhoto saves an uploaded file to the filesystem.
func UploadSavePhoto(f io.Reader, name string) (string, string, error) {
	newID := uuid.NewV4().String()
	ext := filepath.Ext(name)
	path := "data/uploads/photos/" + newID + ext

	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return "", "", err
	}

	err = writeFile(path, fileBytes, 0644)
	if err != nil {
		return "", "", err
	}
	err = createThumbnail(path, fileBytes)
	if err != nil {
		return "", "", err
	}
	return newID, path, err
}

func writeFile(p string, f []byte, fmode os.FileMode) error {
	return ioutil.WriteFile(p, f, fmode)
}

// RemoveFile removes a file from the filesystem.
func RemoveFile(f string) error {
	log.Print("Removing item from filesystem: ", f)
	return os.Remove(f)
}
