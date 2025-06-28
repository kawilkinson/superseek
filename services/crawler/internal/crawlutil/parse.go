package crawlutil

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func ParseInt(value string) (int, error) {
	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("unable to parse integer: %w", err)
	}

	return num, nil
}

func ParseTime(value string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC1123, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse timestamp: %w", err)
	}

	return parsedTime, nil
}

func ParseStringsSlice(value string) ([]string, error) {
	var links []string

	err := json.Unmarshal([]byte(value), &links)
	if err != nil {
		return nil, fmt.Errorf("unable to parse json string slice: %w", err)
	}

	return links, nil
}
