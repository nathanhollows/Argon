package models

import (
	"fmt"
	"html/template"

	"github.com/nathanhollows/Argon/internal/helpers"
	"gorm.io/gorm"
)

// Media stores information about all the different types of media
type Media struct {
	gorm.Model
	Title   string
	File    string
	Type    string
	Format  string
	Hash    string
	Caption string
}

// Shortcode will generate the shortcode that the system convert froms markdown.
// This function will turn a media image into a img src set for accessibly display.
func (media *Media) Shortcode() string {
	return fmt.Sprint("[[", media.Type, ":", media.ID, "]]")
}

// ToHTML generates the HTML for a specific media object
// Returns string for use in the parser
func (media *Media) ToHTML() template.HTML {
	if media.Type == "image" {
		imageTemplate := `<figure><img 
		sizes="(max-width: 2000px) 100vw, 2000px" 
		srcset="
		%s 576w,
		%s 1000w",
		%s 2000w",
		src="%s"
		alt="%s">
		<figcaption>%s</figcaption></figure>`

		imgHTML := fmt.Sprintf(imageTemplate,
			media.ImgURL("small"),
			media.ImgURL("medium"),
			media.ImgURL("large"),
			media.ImgURL(""),
			media.Caption,
			media.Caption)

		return template.HTML(imgHTML)
	}
	return template.HTML("<mark>The file does not exist</mark><br>")
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
