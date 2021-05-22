package jphotos

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/prplecake/jphotos/app"
	"github.com/prplecake/jphotos/auth"
	"github.com/prplecake/jphotos/db"
)

var (
	AlbumsEndpoint = "/albums"
)

func verifyAlbumInput(name string) []string {
	issues := []string{}
	if len(name) == 0 {
		issues = append(issues,
			fmt.Sprintf("Album name cannot be blank."))
	}
	return issues
}

func (s *Server) handleAlbumIndex(w http.ResponseWriter, r *http.Request) {
	type albumData struct {
		Title   string
		Albums  []db.Album
		Auth    *auth.Authorization
		Errors  []string
		Version string
		Branch  string
	}

	auth, _ := auth.Get(r, auth.RoleUser, s.db)

	var errors []string

	if r.Method == "POST" {
		name := r.FormValue("name")
		log.Printf("Album name: %s", name)
		errors = verifyAlbumInput(name)
		fmt.Print("Errors: ", errors)
		if len(errors) == 0 {
			if err := s.db.AddAlbum(name); err != nil {
				if err == db.ErrAlbumExists {
					errors = append(errors,
						fmt.Sprintf("Album name already exists."))
				} else if err == db.ErrAlbumNameInvalid {
					errors = append(errors,
						fmt.Sprintf("Album name is invalid."))
				} else {
					log.Fatal(err)
				}
			}
		}
	}

	albums, err := s.db.GetAllAlbums()
	if err != nil {
		log.Print(err)
		app.ThrowInternalServerError(w)
		return
	}

	app.RenderTemplate(w, "albums", albumData{
		Title:   "Albums",
		Albums:  albums,
		Auth:    auth,
		Errors:  errors,
		Version: app.CurrentVersion,
		Branch:  app.CurrentBranch,
	})

	return
}

type albumData struct {
	Album  *db.Album
	Photos []db.Photo
	Auth   *auth.Authorization
}

type payload struct {
	Title           string
	Album           *db.Album
	Auth            *auth.Authorization
	Photos          []db.Photo
	Version, Branch string
}

func (s *Server) handleGetAlbum(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	log.Print(v)
	album, err := s.db.GetAlbumBySlug(v["slug"])
	if err != nil {
		log.Print(err)
		if err == db.ErrNotFound {
			app.RenderTemplate(w, "error", &app.ErrorInfo{
				Info:          "Album Not Found",
				RedirectLink:  "/",
				RedirectTimer: 3,
			})
		} else {
			app.ThrowInternalServerError(w)
		}
		return
	}

	photos, err := s.db.GetAlbumPhotosByID(album.ID)
	if err != nil {
		log.Print(err)
		app.ThrowInternalServerError(w)
		return
	}

	auth, _ := auth.Get(r, auth.RoleUser, s.db)

	p := &payload{
		Title:   album.Name,
		Album:   album,
		Auth:    auth,
		Photos:  photos,
		Version: app.CurrentVersion,
		Branch:  app.CurrentBranch,
	}
	app.RenderTemplate(w, "album", p)
	return
}

func (s *Server) handleManageAlbumBySlug(w http.ResponseWriter, r *http.Request) {
	auth, _ := auth.Get(r, auth.RoleUser, s.db)
	if auth == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	v := mux.Vars(r)

	album, err := s.db.GetAlbumBySlug(v["slug"])
	if err != nil {
		if err == db.ErrNotFound {
			app.RenderTemplate(w, "error", &app.ErrorInfo{
				Info:          "Album not found",
				RedirectLink:  AlbumsEndpoint,
				RedirectTimer: 3,
			})
		}
		log.Print(err)
		app.ThrowInternalServerError(w)
		return
	}

	//photos, err := s.db.GetAlbumPhotosByID(album.ID)
	//if err != nil {
	//	log.Print(err)
	//	app.ThrowInternalServerError(w)
	//	return
	//}
	if strings.HasSuffix(r.URL.String(), "rename") {
		if len(r.FormValue("new_name")) > 0 {
			err := s.db.RenameAlbumByID(album.ID, r.FormValue("new_name"))
			if err != nil {
				log.Print(err)
				app.ThrowInternalServerError(w)
				return
			}
			slug, err := s.db.GetAlbumSlugByID(album.ID)
			if err != nil {
				log.Print(err)
				app.ThrowInternalServerError(w)
				return
			}
			http.Redirect(w, r, "/album/"+slug, http.StatusSeeOther)
		}
	} else if strings.HasSuffix(r.URL.String(), "update") {
		r.ParseForm()
		log.Print("Update...")
		for k, v := range r.Form {
			if strings.HasPrefix(k, "photo_") {
				// handle delete photo
				if v[0] == "on" {
					log.Print("Deleting photo...")
					log.Print(k)
					photoID := strings.TrimPrefix(k, "photo_")
					s.db.DeletePhotoByID(photoID)
				}
			}
			if strings.HasPrefix(k, "caption_") {
				photoID := strings.TrimPrefix(k, "caption_")
				s.db.UpdatePhotoCaptionByID(photoID, v[0])
				// handle update captions
			}
		}
	}
	http.Redirect(w, r, "/album/"+v["slug"]+"/manage", http.StatusSeeOther)
}

func (s *Server) handleBulkEditAlbumBySlug(w http.ResponseWriter, r *http.Request) {
	auth, _ := auth.Get(r, auth.RoleUser, s.db)
	if auth == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	v := mux.Vars(r)

	album, err := s.db.GetAlbumBySlug(v["slug"])
	if err != nil {
		if err == db.ErrNotFound {
			app.RenderTemplate(w, "error", &app.ErrorInfo{
				Info:          "Album not found",
				RedirectLink:  AlbumsEndpoint,
				RedirectTimer: 3,
			})
		}
		log.Print(err)
		app.ThrowInternalServerError(w)
		return
	}

	photos, err := s.db.GetAlbumPhotosByID(album.ID)
	if err != nil {
		log.Print(err)
		app.ThrowInternalServerError(w)
		return
	}

	p := &payload{
		Title:  album.Name,
		Album:  album,
		Auth:   auth,
		Photos: photos,
	}

	app.RenderTemplate(w, "album-bulk-edit", p)
}

func (s *Server) handleDeleteAlbumBySlug(w http.ResponseWriter, r *http.Request) {
	auth, _ := auth.Get(r, auth.RoleUser, s.db)
	if auth == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	v := mux.Vars(r)
	album, err := s.db.GetAlbumBySlug(v["slug"])
	if err != nil {
		log.Print(err)
		app.ThrowInternalServerError(w)
		return
	}

	log.Printf("Album %s deleted.", v["slug"])
	photos, err := s.db.GetAlbumPhotosByID(album.ID)
	if err != nil {
		log.Print(err)
		app.ThrowInternalServerError(w)
		return
	}

	for _, photo := range photos {
		log.Printf("Removing photo from filesystem [%s]", photo.ID)
		err := app.RemoveFile(s.config.Uploads.Path + photo.Location)
		if err != nil {
			log.Print(err)
			app.ThrowInternalServerError(w)
			return
		}
		err = app.RemoveFile("data/thumbnails/thumb_" + photo.Location)
		if err != nil {
			log.Print(err)
			app.ThrowInternalServerError(w)
			return
		}
		log.Printf("Removing photo from database [%s]", photo.ID)
		err = s.db.DeletePhotoByID(photo.ID)
		if err != nil {
			log.Print(err)
			app.ThrowInternalServerError(w)
			return
		}
	}

	err = s.db.DeleteAlbumBySlug(v["slug"])
	if err != nil {
		log.Print(err)
		app.ThrowInternalServerError(w)
		return
	}

	http.Redirect(w, r, AlbumsEndpoint, http.StatusSeeOther)
}
