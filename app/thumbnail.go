package app

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"path/filepath"

	"golang.org/x/image/draw"
)

// An Image is an image and it's information.
type Image struct {
	Filename    string
	ContentType string
	Data        []byte
	Size        int
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

func createThumbnail(path string, fb []byte) error {
	i, err := process(path, fb)
	thumbPath := "data/thumbnails/thumb_" + filepath.Base(path)

	dst := createRect(i)
	var buffer bytes.Buffer
	switch i.ContentType {
	case "image/jpeg":
		err := jpeg.Encode(&buffer, dst, nil)
		if err != nil {
			return err
		}
		err = writeFile(thumbPath, buffer.Bytes(), 0644)
		if err != nil {
			return err
		}
	case "image/png":
		err := png.Encode(&buffer, dst)
		if err != nil {
			return err
		}
		err = writeFile(thumbPath, buffer.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func process(path string, fb []byte) (*Image, error) {
	contentType := detectContentType(fb)
	log.Print(contentType)

	_, _, err := image.Decode(bytes.NewReader(fb))
	if err != nil {
		return nil, err
	}

	i := &Image{
		Filename:    filepath.Base(path),
		ContentType: contentType,
		Data:        fb,
		Size:        len(fb),
	}

	return i, nil
}

func createRect(i *Image) *image.RGBA {
	img, _, err := image.Decode(bytes.NewReader(i.Data))
	if err != nil {
		log.Print(err)
	}
	var (
		height = img.Bounds().Max.Y
		width  = img.Bounds().Max.X
		y      = 300
		x      = 300 * width / height
	)
	rect := image.Rect(0, 0, x, y)
	dst := image.NewRGBA(rect)
	scaler := draw.ApproxBiLinear
	scaler.Scale(dst, rect, img, img.Bounds(), draw.Over, nil)
	return dst

}
