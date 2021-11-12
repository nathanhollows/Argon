package models

import "gorm.io/gorm"

// Admin stores logins for the admin panel
type Admin struct {
	gorm.Model
	Username string `gorm:"index:idx_name,unique";json:"username"`
	Password string `json:"password"`
}
