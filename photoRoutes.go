package jphotos

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"git.sr.ht/~mjorgensen/jphotos/app"
)

func (s *Server) handleGetPhotoByID(w http.ResponseWriter, r *http.Request) {
	photo, err := s.db.GetPhotoByID(mux.Vars(r)["id"])
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	app.RenderTemplate(w, "photo", photo)
}
