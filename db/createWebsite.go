package db

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Website struct {
	Name        string   `json:"name"`
	Id          string   `json:"id"`
	IconURL     string   `json:"icon_url"`
	Description string   `json:"description"`
	Tags        pq.StringArray `json:"tags" gorm:"type:text[]"`
	Owner int `json:"owner"`
	Created     int `json:"created"`
	Bumps       int `json:"bumps"`
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
