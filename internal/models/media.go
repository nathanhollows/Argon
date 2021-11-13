package models

import (
	"fmt"

	"github.com/nathanhollows/Argon/internal/helpers"
	"gorm.io/gorm"
)

// Media stores information about all the different types of media
type Media struct {
	gorm.Model
	Title  string
	File   string
	Type   string
	Format string
	Hash   string
}

// GenerateShortcode will generate the shortcode that the system convert froms markdown.
// This function will turn a media image into a img src set for accessibly display.
func (media *Media) GenerateShortcode() string {
	return ""
}

// URL returns the URL for the given media object
func (m Media) URL() string {
	return helpers.URL(fmt.Sprint("public/uploads/", m.Type, "/", m.File))
}

// ImgURL returns the url of the resized img
// Accepts "small", "medium", "large"
func (m Media) ImgURL(size string) string {
	if m.Type != "image" {
		return ""
	}
	switch size {
	case "small":
		return helpers.URL(fmt.Sprint("public/uploads/image/small/", m.File))
	case "medium":
		return helpers.URL(fmt.Sprint("public/uploads/image/medium/", m.File))
	case "large":
		return helpers.URL(fmt.Sprint("public/uploads/image/large/", m.File))
	default:
		return helpers.URL(fmt.Sprint("public/uploads/image/", m.File))
	}
}
