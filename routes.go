package jphotos

import (
	"log"
	"net/http"

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
		app.ThrowInternalServerError(w)
		return
	}

	app.RenderTemplate(w, "home", welcomeData{
		Username: auth.User.Username,
		Groups:   groups,
	})
}
