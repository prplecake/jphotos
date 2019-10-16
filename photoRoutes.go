package jphotos

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"git.sr.ht/~mjorgensen/jphotos/app"
	"git.sr.ht/~mjorgensen/jphotos/auth"
	"git.sr.ht/~mjorgensen/jphotos/db"
)

func (s *Server) handlePhotoByID(w http.ResponseWriter, r *http.Request) {

	type photoData struct {
		Photo     *db.Photo
		Auth      *auth.Authorization
		Albums    []db.Album
		AlbumSlug string
	}

	v := mux.Vars(r)
	switch r.Method {
	case "GET":
		auth, _ := auth.Get(r, auth.RoleUser, s.db)
		photo, err := s.db.GetPhotoByID(v["id"])
		if err != nil {
			if err == db.ErrNotFound {

				log.Print(err)
				app.RenderTemplate(w, "error", &app.ErrorInfo{
					Info:          "Photo not found",
					RedirectLink:  "/",
					RedirectTimer: 3,
				})
				return
			}
		}
		album, err := s.db.GetAlbumIDByPhotoID(v["id"])
		if err != nil {
			log.Print(err)
			http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
			return
		}
		albums, err := s.db.GetAllAlbums()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		app.RenderTemplate(w, "photo", &photoData{photo, auth, albums, album})
	case "POST":
		newCaption := r.FormValue("caption")
		if newCaption != "" {
			err := s.db.UpdatePhotoCaption(v["id"], newCaption)
			if err != nil {
				log.Print(err)
				http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
				return
			}
		}
		newAlbum := r.FormValue("new_album")
		if newAlbum != "" {
			err := s.db.UpdatePhotoAlbum(v["id"], newAlbum)
			if err != nil {
				log.Print(err)
				http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/photo/"+v["id"], http.StatusSeeOther)
	}
}

func (s *Server) handleDeletePhotoByID(w http.ResponseWriter, r *http.Request) {
	auth, _ := auth.Get(r, auth.RoleUser, s.db)
	if auth == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var photo *db.Photo
	v := mux.Vars(r)

	photo, err := s.db.GetPhotoByID(v["id"])
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err = app.RemoveFile("data/uploads/photos/" + photo.Location); err != nil {
		log.Printf("File not found at specified location: %w", err)
	}
	if err = app.RemoveFile("data/thumbnails/thumb_" + photo.Location); err != nil {
		log.Printf("File not found at specified location: %w", err)
	}

	err = s.db.DeletePhotoByID(v["id"])
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Print("Photo removed from database.")
	vars := r.URL.Query()
	log.Print(vars)
	next := vars["next"][0]
	log.Print("v: ", v)
	log.Print("Next: ", next)
	if len(next) == 0 {
		next = "/albums"
	}
	http.Redirect(w, r, next, http.StatusSeeOther)
}
