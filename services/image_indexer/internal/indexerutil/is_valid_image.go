package indexerutil

import (
	"image"
	"log"
	"net/http"
	"strings"
)

func IsValidImage(url string, minWidth, minHeight int) bool {
	var absoluteURL string
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		absoluteURL = "https://" + url
	}

	client := http.Client{
		Timeout: Timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("unable to get image %s: %v\n", absoluteURL, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("non 200 response for image %s: %d\n", absoluteURL, resp.StatusCode)
		return false
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Printf("unable to decode image %s: %v\n", absoluteURL, err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	return width >= minWidth && height >= minHeight
}
