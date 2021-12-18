package admin

import (
	"net/http"

	"github.com/nathanhollows/Argon/internal/handler"
)

// Analytics shows player stats
func Analytics(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	data := make(map[string]interface{})
	data["title"] = "Analytics | Admin"
	data["section"] = "analytics"

	return render(w, data, "analytics/index.html")
}
