package models

import "gorm.io/gorm"

// Media stores information about all the different types of media
type Media struct {
	gorm.Model
}

// GenerateShortcode will generate the shortcode that the system convert froms markdown.
// This function will turn a media image into a img src set for accessibly display.
func (media *Media) GenerateShortcode() string {
	return ""
}
