// Argon is a QR code based education platform
package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"
	"github.com/nathanhollows/Argon/internal/filesystem"
	"github.com/nathanhollows/Argon/internal/handler"
	"github.com/nathanhollows/Argon/internal/handler/admin"
	"github.com/nathanhollows/Argon/internal/handler/public"
	"github.com/nathanhollows/Argon/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var router *chi.Mux
var env handler.Env

func init() {
	router = chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Compress(5))

	var store sessions.Store
	if key, ok := os.LookupEnv("ARGON_SESSION_KEY"); ok {
		store = sessions.NewCookieStore([]byte(key))
	} else {
		panic("env var ARGON_SESSION_KEY must be set")
	}

	db, err := gorm.Open(sqlite.Open("trace.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	env = handler.Env{
		Session: store,
		DB:      *db,
		Data:    make(map[string]interface{}),
	}
}

func main() {
	env.DB.AutoMigrate(
		&models.Poll{},
		&models.Admin{},
		&models.ScanEvent{},
		&models.User{},
		&models.Page{},
		&models.Trail{},
		&models.Gallery{},
		&models.Library{},
	)
	routes()
	fmt.Println(http.ListenAndServe(":8050", router))
}

// Set up the routes needed for the game.
func routes() {
	router.Handle("/", handler.HandlePublic{Env: &env, H: public.Index})

	var pages = []models.Page{}
	env.DB.Model(models.Page{}).Where("system").Find(&pages)
	for _, page := range pages {
		router.Handle(fmt.Sprintf("/%v", page.Code), handler.HandlePublic{Env: &env, H: public.RegularPage})
	}

	router.Handle("/library", handler.HandlePublic{Env: &env, H: public.Library})
	router.Handle("/vote", handler.HandlePublic{Env: &env, H: public.Vote})

	router.Handle("/{code:[A-z]{5}}", handler.HandlePublic{Env: &env, H: public.Page})
	router.Handle("/s/{code:[A-z]{5}}", handler.HandlePublic{Env: &env, H: public.Scan})
	router.Handle("/qr/{code:[A-z]{5}} - {fluff}.{format:[A-z]{3}}", handler.HandlePublic{Env: &env, H: public.QR})

	router.Handle("/login", handler.HandlePublic{Env: &env, H: public.Login})
	router.Handle("/logout", handler.HandlePublic{Env: &env, H: public.Logout})

	router.Handle("/admin", handler.HandleAdmin{Env: &env, H: admin.Pages})
	router.Handle("/admin/media", handler.HandleAdmin{Env: &env, H: admin.Media})
	router.Handle("/admin/media/delete/{id}", handler.HandleAdmin{Env: &env, H: admin.DeleteMedia})
	router.Handle("/admin/upload", handler.HandleAdmin{Env: &env, H: admin.Upload})
	router.Handle("/admin/analytics", handler.HandleAdmin{Env: &env, H: admin.Analytics})
	router.Handle("/admin/pages", handler.HandleAdmin{Env: &env, H: admin.Pages})
	router.Handle("/admin/pages/delete", handler.HandleAdmin{Env: &env, H: admin.DeletePage})
	router.Handle("/admin/pages/restore", handler.HandleAdmin{Env: &env, H: admin.Restore})
	router.Handle("/admin/pages/edit/{code}", handler.HandleAdmin{Env: &env, H: admin.EditPage})
	router.Handle("/admin/pages/create", handler.HandleAdmin{Env: &env, H: admin.CreatePage})
	router.Handle("/admin/pages/preview", handler.HandleAdmin{Env: &env, H: admin.PreviewMD})
	router.Handle("/admin/pages/qr/{code}.svg", handler.HandleAdmin{Env: &env, H: admin.QR})

	router.Handle("/404", handler.HandlePublic{Env: &env, H: public.Error404})
	router.NotFound(public.NotFound)

	workDir, _ := os.Getwd()
	filesDir := filesystem.Myfs{Dir: http.Dir(filepath.Join(workDir, "web/static"))}
	filesystem.FileServer(router, "/public", filesDir)
}
