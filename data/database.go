package data

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Hashtag struct {
	ID     uint   `gorm:"primary_key",json:"id"`
	Name   string `json:"name"`
	FeedID int    `json:"feed_id"`
}

var db *gorm.DB
var err error

func init() {
	db, err = gorm.Open("sqlite3", "./database.db")
	if err != nil {
		panic("failed to connect database")
	}
}

func RunMigration() (err error) {
	db.AutoMigrate(&Hashtag{})

	return
}
