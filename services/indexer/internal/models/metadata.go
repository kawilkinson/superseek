package models

import (
	"fmt"
	"log"
	"time"
)

type Metadata struct {
	ID          string
	Title       string
	Description string
	SummaryText string
	LastCrawled time.Time
	Keywords    map[string]int
}

func FromMap(data map[string]interface{}) (*Metadata, error) {
	if data == nil {
		log.Println("no data found for fromMap function for Metadata, skipping...")
		return nil, nil
	}

	lastCrawledStr, isString := data["last_crawled"].(string)
	if !isString {
		return nil, fmt.Errorf("'last_crawled' in metadata must be a string")
	}

	lastCrawled, err := time.Parse(time.RFC1123, lastCrawledStr)
	if err != nil {
		return nil, fmt.Errorf("invalid time format, not able to detect RFC1123: %v", err)
	}

	keywords := make(map[string]int)
	if kwRaw, exists := data["keywords"]; exists {
		if kwMap, isMap := kwRaw.(map[string]interface{}); isMap {
			for key, value := range kwMap {
				if intVal, isFloat64 := value.(float64); isFloat64 {
					keywords[key] = int(intVal)
				}
			}
		}
	}

	return &Metadata{
		ID:          getString(data, "id"),
		Title:       getString(data, "title"),
		Description: getString(data, "description"),
		SummaryText: getString(data, "summary_text"),
		LastCrawled: lastCrawled,
		Keywords:    keywords,
	}, nil
}

// helper function
func getString(data map[string]interface{}, key string) string {
	if val, exists := data[key]; exists {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}

	return ""
}
