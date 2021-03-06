package jphotos

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/prplecake/jphotos/app"
	"github.com/prplecake/jphotos/auth"
	"github.com/prplecake/jphotos/db"
)

func (s *Server) handlePhotoByID(w http.ResponseWriter, r *http.Request) {

	type photoData struct {
		Photo     *db.Photo
		Auth      *auth.Authorization
		Albums    []db.Album
		AlbumSlug string
		Version   string
		Branch    string
		Previous  string
		Next      string
	}

	v := mux.Vars(r)
	switch r.Method {
	case "GET":
		auth, _ := auth.Get(r, auth.RoleUser, s.db)
		photo, err := s.db.GetPhotoByUUID(v["id"])
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
		albumUUID, err := s.db.GetAlbumUUIDByPhotoUUID(v["id"])
		if err != nil {
			log.Print(err)
			app.ThrowInternalServerError(w)
			return
		}
		albumSlug, err := s.db.GetAlbumSlugByUUID(albumUUID)
		if err != nil {
			log.Print(err)
			app.ThrowInternalServerError(w)
			return
		}
		albums, err := s.db.GetAllAlbums()
		if err != nil {
			log.Print(err)
			app.ThrowInternalServerError(w)
			return
		}

		previous := s.db.GetPreviousAlbumPhoto(albumUUID, photo.ID)
		next := s.db.GetNextAlbumPhoto(albumUUID, photo.ID)

		version := app.CurrentVersion
		branch := app.CurrentBranch
		app.RenderTemplate(w, "photo", &photoData{photo, auth, albums, albumSlug, version, branch, previous, next})
	case "POST":
		newCaption := r.FormValue("caption")
		if newCaption != "" {
			err := s.db.UpdatePhotoCaptionByUUID(v["id"], newCaption)
			if err != nil {
				log.Print(err)
				app.ThrowInternalServerError(w)
				return
			}
		}
		newAlbum := r.FormValue("new_album")
		if newAlbum != "" {
			err := s.db.UpdatePhotoAlbum(v["id"], newAlbum)
			if err != nil {
				log.Print(err)
				app.ThrowInternalServerError(w)
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

	photo, err := s.db.GetPhotoByUUID(v["id"])
	if err != nil {
		log.Print(err)
		app.ThrowInternalServerError(w)
		return
	}
	if err = app.RemoveFile("data/uploads/photos/" + photo.Location); err != nil {
		log.Printf("File not found at specified location: %v", err)
	}
	if err = app.RemoveFile("data/thumbnails/thumb_" + photo.Location); err != nil {
		log.Printf("File not found at specified location: %v", err)
	}

	err = s.db.DeletePhotoByUUID(v["id"])
	if err != nil {
		log.Print(err)
		app.ThrowInternalServerError(w)
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
