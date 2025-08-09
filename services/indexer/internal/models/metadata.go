package models

import (
	"log"
)

type Metadata struct {
	ID          string
	Title       string
	Description string
	SummaryText string
	LastCrawled string
	Keywords    map[string]int
}

func FromMap(data map[string]interface{}, pageData *Page) (*Metadata, error) {
	if data == nil {
		log.Println("no data found for fromMap function for Metadata, skipping...")
		return nil, nil
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
		LastCrawled: pageData.LastCrawled.String(),
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
