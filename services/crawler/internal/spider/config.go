package spider

import (
	"fmt"
	"net/url"
	"sync"
)

type Config struct {
	maxPages           int
	Pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	Wg                 *sync.WaitGroup
	runOnce            sync.Once
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
