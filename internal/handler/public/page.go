package public

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/nathanhollows/Argon/internal/flash"
	"github.com/nathanhollows/Argon/internal/handler"
	"github.com/nathanhollows/Argon/internal/helpers"
	"github.com/nathanhollows/Argon/internal/models"
	"gorm.io/gorm/clause"
)

// Page delivers the page relating to a particular code.
// This function does not track scan events.
func Page(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	data := make(map[string]interface{})
	data["section"] = "library"

	code := chi.URLParam(r, "code")
	page := models.Page{}
	env.DB.Where("Code = ?", code).Preload(clause.Associations).Find(&page)
	if page.Code == "" {
		flash.Set(w, r, flash.Message{Message: "That's not a valid code"})
		http.Redirect(w, r, "/404", http.StatusFound)
		return nil
	}

	if !page.Published {
		flash.Set(w, r, flash.Message{Message: "This page is not yet public", Style: "warning"})
	}

	var count int64
	env.DB.Model(models.Page{}).Where("published = true AND gallery_id = ?", page.GalleryID).Count(&count)
	data["count"] = count

	admin, err := env.Session.Get(r, "admin")
	if err == nil || admin.Values["id"] != nil {
		message := fmt.Sprint(`<b>Admin:</b> <a href="`, helpers.URL("admin/pages/edit/"+string(page.Code)), `">Click to edit page</a>`)
		flash.Set(w, r, flash.Message{Message: message})
	}

	session, err := env.Session.Get(r, "uid")
	if err != nil || session.Values["id"] == nil {
		fmt.Println(err)
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

	data["title"] = page.Title
	data["md"] = parseMD(page.Text, &env.DB)
	data["page"] = page

	data["messages"] = flash.Get(w, r)
	return render(w, data, "page/discovered.html")
}

// Scan handles the scanned URL.
// This functions track scan events
func Scan(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	code := strings.ToUpper(chi.URLParam(r, "code"))
	page := models.Page{}
	env.DB.Where("Code = ?", code).Find(&page)
	if page.Code == "" {
		flash.Set(w, r, flash.Message{Message: "That's not a valid code"})
		http.Redirect(w, r, "/404", http.StatusFound)
		return nil
	}

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

	scan := models.ScanEvent{}
	scan.Page = page
	scan.UserID = fmt.Sprint(session.Values["id"])
	env.DB.Model(&models.ScanEvent{}).Create(&scan)

	http.Redirect(w, r, fmt.Sprintf("/%s", page.Code), http.StatusTemporaryRedirect)
	return nil
}
