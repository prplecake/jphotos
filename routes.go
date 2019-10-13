package jphotos

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"

	"git.sr.ht/~mjorgensen/jphotos/app"
	"git.sr.ht/~mjorgensen/jphotos/auth"
	"git.sr.ht/~mjorgensen/jphotos/db"
)

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	type welcomeData struct {
		Username string
		Groups   []db.Group
	}

	auth, _ := auth.Get(r, auth.RoleUser, s.db)
	if auth == nil {
		app.RenderTemplate(w, "landing", nil)
		return
	}

	groups, err := s.db.GetGroupsForUser(auth.User)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	app.RenderTemplate(w, "home", welcomeData{
		Username: auth.User.Username,
		Groups:   groups,
	})
}

type loginData struct {
	Username, password, Next, Error string
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := auth.Get(r, auth.RoleUser, s.db); err == nil {
		// already logged in - redirect
		app.RenderTemplate(w, "error", &app.ErrorInfo{
			Info:          "Already logged in",
			RedirectLink:  "/",
			RedirectTimer: 0,
		})
		return
	}
	var ld = loginData{
		Username: r.FormValue("username"),
		password: r.FormValue("password"),
		Next:     r.FormValue("next"),
	}
	if r.Method == "GET" {
		log.Print("GET Request")
		app.RenderTemplate(w, "login", ld)
	} else {
		log.Print("POST request")
		token, err := auth.NewSession(ld.Username, ld.password, s.db)
		if err != nil {
			if err == auth.ErrInvalidUsernameOrPassword {
				ld.Error = "Invalid username or password"
				log.Println("Error: ", err)
				app.RenderTemplate(w, "login", ld)
				return
			}
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    auth.SessionCookieName,
			Value:   token.Session,
			Expires: token.Expires,
		})
		http.Redirect(w, r, ld.Next, http.StatusSeeOther)
		log.Println("User '" + ld.Username + "' logged in")
		return
	}
}

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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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

func (s *Server) handleUploadPhoto(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("photoFile")
	if err != nil {
		log.Print("Error retrieving the file")
		log.Print(err)
		return
	}
	defer file.Close()
	log.Printf("Uploaded file: %+v\n", handler.Filename)
	log.Printf("File size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)
	log.Printf("Album: %s\n", r.FormValue("album-id"))

	newID := uuid.NewV4().String()
	path := "uploads/photos/" + newID + ".jpeg"

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Print(err)
	}

	err = ioutil.WriteFile(path, fileBytes, 0644)
	if err != nil {
		log.Print(err)
	}

	log.Print("Successfully uploaded file.")
	w.Write([]byte("Successfully uploaded file."))

	return
}

func (s *Server) handleGetPhotoByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not implemented."))
	return
}
