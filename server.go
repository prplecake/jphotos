package jphotos

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"git.sr.ht/~mjorgensen/jphotos/db"
)

// A Server handles routing and dependency injection into the routes.
type Server struct {
	db     db.Store
	router *mux.Router
}

// NewServer creates a Server backed by a backing database
func NewServer(db db.Store) *Server {
	s := &Server{
		db:     db,
		router: mux.NewRouter(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("Request URL: ", r.URL)
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	r := s.router
	r.HandleFunc("/", s.homeHandler)
	r.HandleFunc("/login", s.loginHandler).Methods("GET", "POST")

	r.HandleFunc("/album", s.handleGetAlbum)
	r.HandleFunc("/album/", s.handleGetAlbum)
	r.HandleFunc("/albums", s.handleAlbumIndex)
	r.HandleFunc("/album/{slug}", s.handleGetAlbum)
	r.HandleFunc("/album/{slug}/manage", s.handleManageAlbumBySlug).
		Methods("GET", "POST")
	r.HandleFunc("/album/{slug}/delete", s.handleDeleteAlbumBySlug).
		Methods("POST")

	r.HandleFunc("/photo/{id}", s.handleGetPhotoByID)
	r.HandleFunc("/photo/{id}/manage", s.handleManagePhotoByID)
	r.HandleFunc("/photo/{id}/delete", s.handleDeletePhotoByID).
		Methods("POST")

	r.HandleFunc("/upload", s.handleUploadPhoto).
		Methods("GET", "POST")

	r.PathPrefix("/uploads/photos/").Handler(
		http.StripPrefix("/uploads/photos/",
			http.FileServer(http.Dir("uploads/photos/"))))

	http.Handle("/", r)
}
