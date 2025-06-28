package main

import (
	"flag"
	"log"
	"os"

	"github.com/kawilkinson/search-engine/internal/controllers"
	"github.com/kawilkinson/search-engine/internal/crawlutil"
	"github.com/kawilkinson/search-engine/internal/pages"
	"github.com/kawilkinson/search-engine/internal/redisdb"
	"github.com/kawilkinson/search-engine/internal/spider"
)

func main() {
	maxConcurrency := flag.Int("max-concurrency", 10, "Max number of concurrent channels")
	maxPages := flag.Int("max-pages", 100, "Max number of pages able to be done at once")
	flag.Parse()

	redisHost := loadEnv("REDIS_HOST", "localhost")
	redisPort := loadEnv("REDIS_PORT", "6379")
	redisPassword := loadEnv("REDIS_PASSWORD", "")
	redisDB := loadEnv("REDIS_DB", "0")
	startingURL := loadEnv("STARTING_URL", "https://en.wikipedia.org/wiki/Japan")
	baseURL := "https://en.wikipedia.org"

	db := &redisdb.RedisDatabase{}
	err := db.ConnectToRedis(redisHost, redisPort, redisPassword, redisDB)
	if err != nil {
		log.Printf("unable to connect to redis database: %v\n", err)
	}

	db.PushURLToQueue(startingURL, 0)
	log.Printf("starting queue with %v\n", startingURL)

	pageController := controllers.CreatePageController(db)
	linksController := controllers.CreateLinksController(db)
	imageController := controllers.CreateImageController(db)

	cfg, err := spider.Configure(baseURL, *maxConcurrency, *maxPages)
	if err != nil {
		log.Fatalf("error configuring crawler: %v", err)
	}

	// Main loop for the crawler to continuously run through URLs and push to Redis
	for {
		log.Printf("Checking number of entries in queue...")
		queueSize, err := db.GetIndexerQueueSize()
		if err != nil {
			log.Printf("unable to get indexer queue size: %v\n", err)
			return
		}

		if queueSize >= crawlutil.MaxIndexerQueueSize {
			log.Printf("indexer queue is full, waiting...\n")
			for {
				signal, err := db.PopSignalQueue()
				if err != nil {
					log.Printf("unable to get signal from queue: %v\n", err)
					return
				}

				if signal == crawlutil.ResumeCrawl {
					log.Println("resuming crawl of web pages")
					break
				}
			}
		}

		log.Println("creating workers...")
		cfg.Wg.Add(1)
		go cfg.CrawlPage(db)
		cfg.Wg.Wait()

		pageController.SavePages(cfg)
		linksController.SaveLinks(cfg)
		imageController.SaveImages(cfg)

		cfg.Pages = make(map[string]*pages.Page)
		cfg.Outlinks = make(map[string]*pages.PageNode)
		cfg.Backlinks = make(map[string]*pages.PageNode)
		cfg.Images = make(map[string][]*pages.Image)
	}

}

func loadEnv(key, fallback string) string {
	if envVariable, exists := os.LookupEnv(key); exists {
		return envVariable
	}

	return fallback
}
