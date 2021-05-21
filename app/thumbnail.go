package app

import (
	"log"
	"os"

	"github.com/prplecake/go-thumbnail"
)

var config = thumbnail.Generator{
	Scaler: "CatmullRom",
}

func genThumbnail(path string, fileBytes []byte, contentType string) error {
	log.Println(config)
	gen := thumbnail.NewGenerator(config)

	i, err := gen.NewImageFromByteArray(fileBytes)
	if err != nil {
		return err
	}
	thumbBytes, err := gen.CreateThumbnail(i)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, thumbBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
