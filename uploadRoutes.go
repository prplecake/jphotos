package jphotos

import (
	"log"
	"net/http"
	"path/filepath"

	"git.sr.ht/~mjorgensen/jphotos/app"
	"git.sr.ht/~mjorgensen/jphotos/auth"
	"git.sr.ht/~mjorgensen/jphotos/db"
)

func (s *Server) handleUploadPhoto(w http.ResponseWriter, r *http.Request) {
	auth, _ := auth.Get(r, auth.RoleUser, s.db)
	if auth == nil {
		log.Print("Error: Unauthorized")
		http.Error(w, "Unauthorized.", 401)
		return
	}
	err := r.ParseMultipartForm(100000)
	if err != nil {
		log.Print(err)
		app.ThrowInternalServerError(w)
		return
	}

	count := 0
	files := r.MultipartForm.File["photoFiles"]
	log.Printf("Preparing to upload %d files.", len(files))
	for i := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			log.Print("Error retrieving the file")
			log.Print(err)
			return
		}
		if files[i].Size == 0 {
			log.Print("Error: file is empty. File size: ", files[i].Size)
			continue
		}

		log.Printf("Uploaded file: %+v\n", files[i].Filename)
		log.Printf("File size: %+v\n", files[i].Size)
		log.Printf("MIME Header: %+v\n", files[i].Header)
		log.Printf("Album: %s\n", r.FormValue("album-id"))

		newID, path, err := app.UploadSavePhoto(file, files[i].Filename, s.config.Uploads)
		if err != nil {
			if err == app.ErrBadContentType {
				log.Printf("Bad content type, not uploading [%s]", files[i].Filename)
				continue
			}
			log.Print("UploadSavePhoto: ", err)
		}

		err = s.db.AddPhoto(db.Photo{
			ID:       newID,
			Caption:  r.FormValue("caption"),
			Location: filepath.Base(path),
		}, r.FormValue("album-id"))

		count++
		log.Printf("Successfully uploaded file %d.", count)
	}
	log.Printf("Uploaded %d of %d files.", count, len(files))
	log.Printf("Uploaded %d files", count)

	slug, err := s.db.GetAlbumSlugByID(r.FormValue("album-id"))
	log.Print(slug)
	if err != nil {
		log.Print(err)
		app.RenderTemplate(w, "error", &app.ErrorInfo{
			Info:          "Album Not Found",
			RedirectLink:  "/",
			RedirectTimer: 3,
		})
	}
	http.Redirect(w, r, "/album/"+slug, http.StatusSeeOther)
}
