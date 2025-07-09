package models

import "log"

type Image struct {
	ID       string
	PageURL  string
	Alt      string
	Filename string
}

func FromHash(image map[string]string, imageURL string) *Image {
	if image == nil {
		log.Println("unable to get image data for FromHash, no data found")
		return nil
	}

	return &Image{
		ID:       imageURL,
		PageURL:  getStringFromMap(image, "page_url"),
		Alt:      getStringFromMap(image, "alt"),
		Filename: getStringFromMap(image, "filename"),
	}
}

func getStringFromMap(currMap map[string]string, key string) string {
	if val, exists := currMap[key]; exists {
		return val
	}

	return ""
}

func (img *Image) ToMap() map[string]string {
	return map[string]string{
		"_id":      img.ID,
		"page_url": img.PageURL,
		"alt":      img.Alt,
		"filename": img.Filename,
	}
}
