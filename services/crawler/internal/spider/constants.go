package spider

import "time"

const (
	Timeout = 4 * time.Second

	IndexerQueueKey     = "indexer_queue"
	SignalQueueKey      = "signal_queue"
	CrawlerQueueKey     = "crawler_queue"
	NormalizedURLPrefix = "normalized_url"
)
