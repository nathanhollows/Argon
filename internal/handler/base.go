package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/nathanhollows/Argon/internal/helpers"
	"gorm.io/gorm"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Status returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

// HandlePublic takes both a game manager and http function.
type HandlePublic struct {
	Env *Env
	H   func(e *Env, w http.ResponseWriter, r *http.Request) error
}

// HandleAdmin takes both a game manager and http function.
type HandleAdmin struct {
	Env *Env
	H   func(e *Env, w http.ResponseWriter, r *http.Request) error
}

// Env is the shared game manager for each request.
type Env struct {
	Session sessions.Store
	DB      gorm.DB
	Data    map[string]interface{}
}

func (h HandlePublic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.Env, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

func (h HandleAdmin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := h.Env.Session.Get(r, "admin")
	if err != nil || session.Values["id"] == nil {
		http.Redirect(w, r, helpers.URL("login"), 302)
	}

	err = h.H(h.Env, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}
