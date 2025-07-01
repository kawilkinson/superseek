package indexerutil

import "time"

const (
	Timeout               = 4 * time.Second
	MaxIndexerQueueSize   = 5000
	MaxScore              = 10000
	MinScore              = -1000
	ImgCtrlExpirationTime = 1

	IndexerQueueKey = "indexer_queue"
	SignalQueueKey  = "signal_queue"
	CrawlerQueueKey = "crawler_queue"
	ResumeCrawl     = "RESUME_CRAWL"

	NormalizedURLPrefix = "normalized_url"
	PagePrefix          = "page_data"
	ImagePrefix         = "image_data"
	PageImagesPrefix    = "page_images"
	BacklinksPrefix     = "backlinks"
	OutlinksPrefix      = "outlinks"
)
