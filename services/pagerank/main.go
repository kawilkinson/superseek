package main

import "os"

func loadEven(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}

// entry into PageRank service using the PageRank algorithm
func main() {

}
