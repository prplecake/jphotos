package jphotos

import (
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
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	r := s.router
	r.HandleFunc("/", s.homeHandler)
	r.HandleFunc("/login", s.loginHandler).Methods("GET", "POST")
	r.HandleFunc("/albums", s.handleAlbumIndex)
	r.HandleFunc("/album/{slug}", s.handleGetAlbum)
	r.HandleFunc("/photo/{id}", s.handleGetPhotoByID)
	r.HandleFunc("/upload", s.handleUploadPhoto).Methods("POST")

	http.Handle("/", r)
}
