package spider

import (
	"fmt"
	"net/url"
	"sync"
)

type Config struct {
	mu                 *sync.Mutex
	Wg                 *sync.WaitGroup
	runOnce            sync.Once
	Pages              map[string]int
	Outlinks           map[string]int
	Backlinks          map[string]int
	Images             map[string][]int
	baseURL            *url.URL
	concurrencyControl chan struct{}
	maxPages           int
}

func (cfg *Config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	_, visited := cfg.Pages[normalizedURL]
	if visited {
		cfg.Pages[normalizedURL] += 1
		return false
	}

	cfg.Pages[normalizedURL] = 1
	return true
}

func (cfg *Config) isMaxPagesReached() (maxReached bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	return len(cfg.Pages) >= cfg.maxPages
}

func Configure(baseURL string, maxConcurrency, maxPages int) (*Config, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing provided base URL: %v", err)
	}

	return &Config{
		maxPages:           maxPages,
		Pages:              make(map[string]int),
		baseURL:            parsedBaseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		Wg:                 &sync.WaitGroup{},
	}, nil
}
