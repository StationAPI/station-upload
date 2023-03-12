package db

import (
	"gorm.io/gorm"
)

type Website struct {
	Name string `json:"name"`
	Id string `json:"id"`
	IconURL string `json:"icon_url"`
	Description string `json:"description"`
	Tags []string `json:"tags"` 
}

func CreateWebsite(website Website, db gorm.DB) {
	db.Create(website)
}

func GetWebsite(id string, db gorm.DB) (Website, bool) {
	website := Website{}

	db.Where("id = ?", id).First(&website)

	if website.Id == "" {
		return website, false
	}

	return website, true
}
