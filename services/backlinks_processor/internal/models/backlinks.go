package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Backlinks struct {
	ID    string              `bson:"_id"`
	Links map[string]struct{} `bson:"links"`
}

func (b *Backlinks) ToMap() bson.M {
	links := make([]string, 0, len(b.Links))
	for link := range b.Links {
		links = append(links, link)
	}

	return bson.M{
		"_id":   b.ID,
		"links": links,
	}
}
