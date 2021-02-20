package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Start begins the game for the team. Prints out their first clue
func Start(env *Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	type Data struct {
		Team string
	}

	r.ParseForm()
	team := r.PostForm.Get("code")
	data := Data{
		Team: team,
	}

	var page string
	if env.Manager.CheckTeam(team) {
		fmt.Println("Team ", team, " has checked in")
		page = "../web/template/start/index.html"
	} else {
		page = "../web/template/index/error.html"
	}

	templates := template.Must(template.ParseFiles(
		"../web/template/index.html",
		page,
	))

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return nil
}
