package jphotos

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"git.sr.ht/~mjorgensen/jphotos/app"
	"git.sr.ht/~mjorgensen/jphotos/auth"
	"git.sr.ht/~mjorgensen/jphotos/db"
)

func (s *Server) handleGetPhotoByID(w http.ResponseWriter, r *http.Request) {

	type photoData struct {
		Photo *db.Photo
		Auth  *auth.Authorization
	}
	auth, _ := auth.Get(r, auth.RoleUser, s.db)
	photo, err := s.db.GetPhotoByID(mux.Vars(r)["id"])
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	app.RenderTemplate(w, "photo", photoData{photo, auth})
}

func (s *Server) handleManagePhotoByID(w http.ResponseWriter, r *http.Request) {
	type managePhotoData struct {
		Photo  *db.Photo
		Auth   *auth.Authorization
		Albums []db.Album
	}

	switch r.Method {
	case "GET":
		auth, _ := auth.Get(r, auth.RoleUser, s.db)
		photo, err := s.db.GetPhotoByID(mux.Vars(r)["id"])
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		albums, err := s.db.GetAlbums()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		app.RenderTemplate(w, "photo_manage", &managePhotoData{photo, auth, albums})
	case "POST":
		log.Print(r)
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
	if err = os.Remove(photo.Location); err != nil {
		log.Printf("File not found at specified location: %-w", err)
	}

	err = s.db.DeletePhotoByID(v["id"])
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Print("Photo removed from database.")

	http.Redirect(w, r, "/albums", http.StatusSeeOther)
}
