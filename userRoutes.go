package jphotos

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/prplecake/jphotos/app"
	"github.com/prplecake/jphotos/auth"
	"github.com/prplecake/jphotos/db"
)

type loginData struct {
	Username, password string
	Next               string
	Error              string
	Version, Branch    string
	Title              string
	Auth               *auth.Authorization
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
		Username: strings.ToLower(r.FormValue("username")),
		password: r.FormValue("password"),
		Next:     r.FormValue("next"),
		Title:    "Login",
		Version:  app.CurrentVersion,
		Branch:   app.CurrentBranch,
	}
	switch r.Method {
	case "GET":
		log.Print("GET Request")
		app.RenderTemplate(w, "login", ld)

	case "POST":
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

type usersData struct {
	Title           string
	Auth            *auth.Authorization
	Users           []db.User
	User            db.User
	Version, Branch string
}

func (s *Server) handleUsersIndex(w http.ResponseWriter, r *http.Request) {
	auth, err := auth.Get(r, auth.RoleUser, s.db)
	if err != nil {
		// not logged in, redirect
		app.RenderTemplate(w, "error", &app.ErrorInfo{
			Info:          "Unauthorized.",
			RedirectLink:  "/",
			RedirectTimer: 0,
		})
	}

	users, err := s.db.GetAllUsers()
	if err != nil {
		log.Print(err)
		app.ThrowInternalServerError(w)
	}
	app.RenderTemplate(w, "users", usersData{
		Title:   "Manage Users",
		Auth:    auth,
		Users:   users,
		Version: app.CurrentVersion,
		Branch:  app.CurrentBranch,
	})

}

func (s *Server) handleGetUserByUsername(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	auth, err := auth.Get(r, auth.RoleUser, s.db)
	if err != nil {
		// not logged in, redirect
		app.RenderTemplate(w, "error", &app.ErrorInfo{
			Info:          "Unauthorized.",
			RedirectLink:  "/",
			RedirectTimer: 0,
		})
	}
	user, err := s.db.GetUserByUsername(v["username"])
	if err != nil {
		log.Print(err)
		app.ThrowInternalServerError(w)
	}
	app.RenderTemplate(w, "user", usersData{
		Title:   "User: " + user.Username,
		Auth:    auth,
		User:    *user,
		Version: app.CurrentVersion,
		Branch:  app.CurrentBranch,
	})
}
