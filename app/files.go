package app

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	uuid "github.com/satori/go.uuid"
)

// UploadSaveFile saves an uploaded file to the filesystem.
func UploadSaveFile(f io.Reader) (string, string, error) {
	newID := uuid.NewV4().String()
	path := "uploads/photos/" + newID + ".jpeg"

	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return "", "", err
	}

	err = ioutil.WriteFile(path, fileBytes, 0644)
	if err != nil {
		return "", "", err
	}
	return newID, path, err
}

// RemoveFile removes a file from the filesystem.
func RemoveFile(f string) error {
	log.Print("Removing item from filesystem: ", f)
	return os.Remove(f)
}
