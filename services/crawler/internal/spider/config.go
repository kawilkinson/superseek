package spider

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
	"github.com/kawilkinson/search-engine/internal/pages"
)

type Config struct {
	mu                 *sync.Mutex
	Wg                 *sync.WaitGroup
	runOnce            sync.Once
	Pages              map[string]*pages.Page
	Outlinks           map[string]*pages.PageNode
	Backlinks          map[string]*pages.PageNode
	Images             map[string][]*pages.Image
	baseURL            *url.URL
	concurrencyControl chan struct{}
	maxPages           int
	MaxConcurrency     int
}

// func (cfg *Config) addPageVisit(normalizedURL string) (isFirst bool) {
// 	cfg.mu.Lock()
// 	defer cfg.mu.Unlock()

// 	_, visited := cfg.Pages[normalizedURL]
// 	if visited {
// 		cfg.Pages[normalizedURL] += 1
// 		return false
// 	}

// 	cfg.Pages[normalizedURL] = 1
// 	return true
// }

func (cfg *Config) isMaxPagesReached() (maxReached bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	return len(cfg.Pages) >= cfg.maxPages
}

func Configure(baseURL string, maxConcurrency *int, maxPages int) (*Config, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing provided base URL: %v", err)
	}

	return &Config{
		mu:                 &sync.Mutex{},
		Wg:                 &sync.WaitGroup{},
		Pages:              make(map[string]*pages.Page),
		Outlinks:           make(map[string]*pages.PageNode),
		Backlinks:          make(map[string]*pages.PageNode),
		Images:             make(map[string][]*pages.Image),
		baseURL:            parsedBaseURL,
		concurrencyControl: make(chan struct{}, *maxConcurrency),
		maxPages:           maxPages,
		MaxConcurrency:     *maxConcurrency,
	}, nil
}

func (cfg *Config) addPage(page *pages.Page) error {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	normalizedURL := page.NormalizedURL

	if _, visited := cfg.Pages[normalizedURL]; visited {
		return fmt.Errorf("page %v already visited", page.NormalizedURL)
	}

	if len(cfg.Pages) >= cfg.maxPages {
		return fmt.Errorf("max pages has been reached")
	}

	cfg.Pages[normalizedURL] = page
	return nil
}

func (cfg *Config) UpdateLinks(normalizedURL string, outgoingLinks []string) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.Outlinks[normalizedURL] = pages.CreatePageNode(normalizedURL)
	for _, link := range outgoingLinks {
		if crawlutil.IsValidURL(link) {
			normalizedOutgoingURL, err := crawlutil.NormalizeURL(link)
			if err != nil {
				continue
			}

			if normalizedOutgoingURL == normalizedURL {
				continue
			}

			if _, exists := cfg.Backlinks[normalizedOutgoingURL]; !exists {
				cfg.Backlinks[normalizedOutgoingURL] = pages.CreatePageNode(normalizedOutgoingURL)
			}

			cfg.Backlinks[normalizedOutgoingURL].AppendLink(normalizedURL)
			cfg.Outlinks[normalizedURL].AppendLink(normalizedOutgoingURL)
		}
	}
}

func (cfg *Config) AddImages(normalizedURL string, images map[string]map[string]string) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	for imgURL, imgAttrs := range images {
		imgAlt := ""
		if alt, exists := imgAttrs["alt"]; exists {
			imgAlt = alt
		}

		image := &pages.Image{
			NormalizedPageURL:   normalizedURL,
			NormalizedSourceURL: imgURL,
			Alt:                 imgAlt,
		}

		cfg.Images[normalizedURL] = append(cfg.Images[normalizedURL], image)
	}
}
