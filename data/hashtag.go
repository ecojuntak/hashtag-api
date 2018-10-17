package data

import (
	"bytes"
	"encoding/json"
	"log"
)

type Feed struct {
	ID      int
	Caption string
}

func filterHastag(caption string) []string {
	var hashtags []string
	buffer := bytes.NewBufferString("")

	for _, c := range caption {
		if buffer.String() != "" {
			if c == ' ' {
				hashtags = append(hashtags, buffer.String())
				buffer = bytes.NewBufferString("")
				continue
			} else if c == '#' {
				hashtags = append(hashtags, buffer.String())
				buffer = bytes.NewBufferString("#")
				continue
			}
			buffer.WriteString(string(c))
		} else {
			if c == '#' {
				buffer.WriteString(string(c))
			}
		}
	}

	hashtags = append(hashtags, buffer.String())
	return hashtags
}

func GetFeedIds(name string) (feed_ids []int) {
	hashtags := []Hashtag{}
	db.Where("name = ?", "#"+name).Find(&hashtags)

	for _, hashtag := range hashtags {
		feed_ids = append(feed_ids, hashtag.FeedID)
	}

	return
}

func ProcessMessage(msg string) {
	feed := &Feed{}
	err := json.Unmarshal([]byte(msg), feed)

	hashtags := filterHastag(feed.Caption)

	for _, hashtag := range hashtags {
		store(&Hashtag{
			FeedID: feed.ID,
			Name:   hashtag,
		})

		log.Println(hashtag + " stored")
	}

	if err != nil {
		log.Println(err)
	}
}

func store(h *Hashtag) {
	db.Create(h)
}

func GetAll() (hashtags []Hashtag) {
	db.Find(&hashtags)
	return hashtags
}
