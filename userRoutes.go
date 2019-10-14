package jphotos

import (
	"log"
	"net/http"
	"strings"

	"git.sr.ht/~mjorgensen/jphotos/app"
	"git.sr.ht/~mjorgensen/jphotos/auth"
)

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
		Username: strings.ToLower(r.FormValue("username")),
		password: r.FormValue("password"),
		Next:     r.FormValue("next"),
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
