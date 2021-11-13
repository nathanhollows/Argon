package public

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nathanhollows/Argon/internal/flash"
	"github.com/nathanhollows/Argon/internal/handler"
	"github.com/nathanhollows/Argon/internal/models"
)

// Index is the homepage of the game.
// Prints a very simple page asking only for a team code.
func Index(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-store")
	data := make(map[string]interface{})
	data["messages"] = flash.Get(w, r)
	data["section"] = "index"

	session, err := env.Session.Get(r, "uid")
	if err != nil || session.Values["id"] == nil {
		session, err = env.Session.New(r, "uid")
		session.Options.HttpOnly = true
		session.Options.SameSite = http.SameSiteStrictMode
		session.Options.Secure = true
		id := uuid.New()
		session.Values["id"] = id.String()
		session.Save(r, w)
	}

	var trails []models.ResultsTrailCounts
	env.DB.Raw(models.QueryTrailCountByUser, session.Values["id"]).Scan(&trails)
	data["trails"] = trails

	var pages []models.ResultFoundPages
	env.DB.Raw(models.QueryFindPagesByUser, session.Values["id"]).Scan(&pages)
	for i, _ := range pages {
		env.DB.Model(&models.Media{}).Where("id = ?", pages[i].CoverID).Limit(1).Find(&pages[i].Cover)
	}
	data["pages"] = pages
	return render(w, data, "index/index.html")
}
