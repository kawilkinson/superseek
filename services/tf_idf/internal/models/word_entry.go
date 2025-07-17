package models

type WordEntry struct {
	URL  string  `bson:"url"`
	TF   float64 `bson:"tf"`
	Word string  `bson:"word,omitempty"`
}
