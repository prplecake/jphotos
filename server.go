package jphotos

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/prplecake/jphotos/app"
	"github.com/prplecake/jphotos/db"
)

// A Server handles routing and dependency injection into the routes.
type Server struct {
	db     db.Store
	config app.Configuration
	router *mux.Router
}

// NewServer creates a Server backed by a backing database
func NewServer(db db.Store, config app.Configuration) *Server {
	s := &Server{
		db:     db,
		config: config,
		router: mux.NewRouter(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// log.Print("Request URL: ", r.URL)
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	r := s.router
	r.HandleFunc("/", s.handleAlbumIndex)
	r.HandleFunc("/login", s.loginHandler).Methods("GET", "POST")

	r.HandleFunc("/album", s.handleGetAlbum)
	r.HandleFunc("/album/", s.handleGetAlbum)
	r.HandleFunc("/albums", s.handleAlbumIndex)
	r.HandleFunc("/album/{slug}", s.handleGetAlbum)
	r.HandleFunc("/album/{slug}/manage", s.handleBulkEditAlbumBySlug).
		Methods("GET")
	r.HandleFunc("/album/{slug}/rename", s.handleManageAlbumBySlug).
		Methods("POST")
	r.HandleFunc("/album/{slug}/update", s.handleManageAlbumBySlug).
		Methods("POST")
	r.HandleFunc("/album/{slug}/delete", s.handleDeleteAlbumBySlug).
		Methods("POST")

	r.HandleFunc("/photo/{id}", s.handlePhotoByID)
	r.HandleFunc("/photo/{id}/manage", s.handlePhotoByID)
	r.HandleFunc("/photo/{id}/delete", s.handleDeletePhotoByID).
		Methods("POST")

	r.HandleFunc("/upload", s.handleUploadPhoto).
		Methods("GET", "POST")

	r.HandleFunc("/users", s.handleUsersIndex)
	r.HandleFunc("/user/{username}", s.handleGetUserByUsername)

	r.PathPrefix("/p/").Handler(
		http.StripPrefix("/p/",
			http.FileServer(http.Dir("data/uploads/photos/"))))
	r.PathPrefix("/t/").Handler(
		http.StripPrefix("/t/",
			http.FileServer(http.Dir("data/thumbnails/"))))
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("assets/"))))

	http.Handle("/", r)
	log.Print("Routes initialized.")
}
