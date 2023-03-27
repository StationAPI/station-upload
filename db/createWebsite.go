package db

import (
  "time"

  "github.com/google/uuid"
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
  LastBumped  int `json:"last_bumped"`
	Bumps       int `json:"bumps"`
  Upvotes int `json:"upvotes"`
  UpvotesToday int `json:"upvotes_today"`
}

func CreateWebsite(website Website, db gorm.DB) {
  website.Id = uuid.NewString()
	website.Bumps = 0
	website.Created = int(time.Now().Unix())
  website.LastBumped = 0
  website.Upvotes = 0
  website.UpvotesToday = 0


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
