package admin

import (
	"html/template"
	"log"
	"net/http"

	"github.com/nathanhollows/AmazingTrace/pkg/handler"
)

// Admin handles the teams and shows the current status
func Admin(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	templates := template.Must(template.ParseFiles(
		"../web/template/index.html",
		"../web/template/admin/index.html"))

	if err := templates.ExecuteTemplate(w, "base", env.Manager); err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return nil
}

// FastForward completes three clues for a team.
func FastForward(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	r.ParseForm()
	teamCode := r.PostFormValue("code")

	index, _ := env.Manager.GetTeam(teamCode)
	team := &env.Manager.Teams[index]
	team.FastForward()
	http.Redirect(w, r, "/admin", 301)
	return nil
}

// Hinder completes three clues for a team.
func Hinder(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	r.ParseForm()
	teamCode := r.PostFormValue("code")

	index, _ := env.Manager.GetTeam(teamCode)
	team := &env.Manager.Teams[index]
	team.Hinder()
	http.Redirect(w, r, "/admin", 301)
	return nil
}

// CreateTeam completes three clues for a team.
func CreateTeam(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	env.Manager.CreateTeams(3)
	http.Redirect(w, r, "/admin", 301)
	return nil
}

// Codes list all available codes
func Codes(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	templates := template.Must(template.ParseFiles(
		"../web/template/index.html",
		"../web/template/admin/codes.html"))

	if err := templates.ExecuteTemplate(w, "base", env.Manager); err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return nil
}