package crawler_utilities

import "time"

const (
	Timeout             = 4 * time.Second
	MaxIndexerQueueSize = 5000

	IndexerQueueKey     = "indexer_queue"
	SignalQueueKey      = "signal_queue"
	CrawlerQueueKey     = "crawler_queue"
	NormalizedURLPrefix = "normalized_url"
	ResumeCrawl         = "RESUME_CRAWL"
)
