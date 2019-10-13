package jphotos

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"git.sr.ht/~mjorgensen/jphotos/app"
	"git.sr.ht/~mjorgensen/jphotos/auth"
	"git.sr.ht/~mjorgensen/jphotos/db"
)

func (s *Server) handleAlbumIndex(w http.ResponseWriter, r *http.Request) {
	type albumData struct {
		Albums []db.Album
		Auth   *auth.Authorization
	}

	if r.Method == "POST" {
		name := r.FormValue("name")
		if err := s.db.AddAlbum(name); err != nil {
			log.Fatal(err)
		}

	}

	auth, _ := auth.Get(r, auth.RoleUser, s.db)

	albums, err := s.db.GetAlbums()
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	app.RenderTemplate(w, "albums", albumData{
		Albums: albums,
		Auth:   auth,
	})

	return
}

func (s *Server) handleGetAlbum(w http.ResponseWriter, r *http.Request) {
	type albumData struct {
		Album  *db.Album
		Photos []db.Photo
		Auth   *auth.Authorization
	}

	album, err := s.db.GetAlbum(mux.Vars(r)["slug"])
	if err != nil {
		log.Print(err)
		if err == db.ErrNotFound {
			app.RenderTemplate(w, "error", &app.ErrorInfo{
				Info:          "Album Not Found",
				RedirectLink:  "/",
				RedirectTimer: 3,
			})
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	photos, err := s.db.GetAlbumPhotos(album.ID)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	auth, _ := auth.Get(r, auth.RoleUser, s.db)

	var ad = albumData{
		Album:  album,
		Photos: photos,
		Auth:   auth,
	}
	app.RenderTemplate(w, "album", ad)
	return
}
