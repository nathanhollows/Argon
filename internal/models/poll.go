package models

import "gorm.io/gorm"

// Poll keeps track of results
type Poll struct {
	gorm.Model
	Poll   uint
	Result string
	UserID string `sql:"DEFAULT:NULL"`
}
