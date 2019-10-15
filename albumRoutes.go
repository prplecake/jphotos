package jphotos

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"git.sr.ht/~mjorgensen/jphotos/app"
	"git.sr.ht/~mjorgensen/jphotos/auth"
	"git.sr.ht/~mjorgensen/jphotos/db"
)

func verifyAlbumInput(name string) []string {
	issues := []string{}
	if len(name) == 0 {
		log.Printf("Album name [%s] verification...", name)
		issues = append(issues,
			fmt.Sprintf("Album name cannot be blank."))
	}
	return issues
}

func (s *Server) handleAlbumIndex(w http.ResponseWriter, r *http.Request) {
	type albumData struct {
		Albums []db.Album
		Auth   *auth.Authorization
		Errors []string
	}

	auth, _ := auth.Get(r, auth.RoleUser, s.db)

	if r.Method == "POST" {
		albums, err := s.db.GetAlbums()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		name := r.FormValue("name")
		log.Printf("Album name: %s", name)
		errors := verifyAlbumInput(name)
		fmt.Print("Errors: ", errors)
		if len(errors) > 0 {
			app.RenderTemplate(w, "albums", albumData{
				Albums: albums,
				Auth:   auth,
				Errors: errors,
			})
			return
		}
		if err := s.db.AddAlbum(name); err != nil {
			log.Fatal(err)
		}
	}

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

type albumData struct {
	Album  *db.Album
	Photos []db.Photo
	Auth   *auth.Authorization
}

type payload struct {
	Title  string
	Album  *db.Album
	Auth   *auth.Authorization
	Photos []db.Photo
}

func (s *Server) handleGetAlbum(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	log.Print(v)
	album, err := s.db.GetAlbum(v["slug"])
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

	p := &payload{
		Title:  album.Name,
		Album:  album,
		Auth:   auth,
		Photos: photos,
	}
	app.RenderTemplate(w, "album", p)
	return
}

func (s *Server) handleManageAlbumBySlug(w http.ResponseWriter, r *http.Request) {
	var manageAlbumData albumData
	auth, _ := auth.Get(r, auth.RoleUser, s.db)
	if auth == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	v := mux.Vars(r)

	album, err := s.db.GetAlbum(v["slug"])
	if err != nil {
		if err == db.ErrNotFound {
			app.RenderTemplate(w, "error", &app.ErrorInfo{
				Info:          "Album not found",
				RedirectLink:  "/albums",
				RedirectTimer: 3,
			})
		}
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	photos, err := s.db.GetAlbumPhotos(album.ID)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "GET":
		manageAlbumData.Album = album
		manageAlbumData.Auth = auth
		manageAlbumData.Photos = photos

		app.RenderTemplate(w, "album_manage", manageAlbumData)
	case "POST":
		if len(r.FormValue("new_name")) > 0 {
			err := s.db.RenameAlbum(album.ID, r.FormValue("new_name"))
			if err != nil {
				log.Print(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			slug, err := s.db.GetAlbumSlugByID(album.ID)
			if err != nil {
				log.Print(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/album/"+slug, http.StatusSeeOther)
		}
	}
}

func (s *Server) handleDeleteAlbumBySlug(w http.ResponseWriter, r *http.Request) {
	auth, _ := auth.Get(r, auth.RoleUser, s.db)
	if auth == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	v := mux.Vars(r)
	album, err := s.db.GetAlbum(v["slug"])
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Album %s deleted.", v["slug"])
	photos, err := s.db.GetAlbumPhotos(album.ID)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	for _, photo := range photos {
		log.Printf("Removing photo from filesystem [%s]", photo.ID)
		err := app.RemoveFile(photo.Location)
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		log.Printf("Removing photo from database [%s]", photo.ID)
		err = s.db.DeletePhotoByID(photo.ID)
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	err = s.db.DeleteAlbumBySlug(v["slug"])
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/albums", http.StatusSeeOther)
}
