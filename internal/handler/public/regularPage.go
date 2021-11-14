package public

import (
	"net/http"
	"strings"

	"github.com/nathanhollows/Argon/internal/flash"
	"github.com/nathanhollows/Argon/internal/handler"
	"github.com/nathanhollows/Argon/internal/helpers"
	"github.com/nathanhollows/Argon/internal/models"
)

// RegularPage renders a boring ol' page without any gameplay mechanics
func RegularPage(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	data := make(map[string]interface{})
	page := models.Page{}
	code := strings.ToLower(
		strings.TrimPrefix(
			strings.TrimSuffix(
				r.URL.String(), "/"), "/"))
	result := env.DB.Where("code = ?", code).Find(&page)
	if result.RowsAffected == 0 {
		http.Redirect(w, r, helpers.URL("404"), 404)
		return nil
	}
	data["md"] = parseMD(page.Text, &env.DB)
	data["title"] = page.Title
	data["messages"] = flash.Get(w, r)
	return render(w, data, "page/default.html")
}
