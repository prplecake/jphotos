package jphotos

import (
	"io/ioutil"
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"

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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["photoFiles"]
	for i := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			log.Print("Error retrieving the file")
			log.Print(err)
			return
		}

		log.Printf("Uploaded file: %+v\n", files[i].Filename)
		log.Printf("File size: %+v\n", files[i].Size)
		log.Printf("MIME Header: %+v\n", files[i].Header)
		log.Printf("Album: %s\n", r.FormValue("album-id"))

		newID := uuid.NewV4().String()
		path := "uploads/photos/" + newID + ".jpeg"

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Print(err)
		}

		err = ioutil.WriteFile(path, fileBytes, 0644)
		if err != nil {
			log.Print(err)
		}

		err = s.db.AddPhoto(db.Photo{
			ID:       newID,
			Caption:  r.FormValue("caption"),
			Location: path,
		}, r.FormValue("album-id"))

		log.Print("Successfully uploaded file.")
	}

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
