package indexerutil

import "time"

const (
	Timeout       = 4 * time.Second
	MaxIndexWords = 1000

	IndexerQueueKey      = "indexer_queue"
	ImageIndexerQueueKey = "image_indexer_queue"
	SignalQueueKey       = "signal_queue"
	ResumeCrawl          = "RESUME_CRAWL"

	NormalizedURLPrefix = "normalized_url"
	UrlMetadataPrefix   = "url_metadata"
	PagePrefix          = "page_data"
	ImagePrefix         = "image_data"
	WordPrefix          = "word"
	PageImagesPrefix    = "page_images"
	WordImagesPrefix    = "word_images"
	BacklinksPrefix     = "backlinks"
	OutlinksPrefix      = "outlinks"
)

var fileTypes = map[string]struct{}{
	"png": {}, "svg": {}, "ico": {}, "gif": {}, "jpeg": {}, "jpg": {},
}

// this is a general collection of domains for the indexer to watch out for
var popularDomains = map[string]struct{}{
	// top level domains and country codes
	"com": {}, "org": {}, "net": {}, "edu": {}, "gov": {}, "mil": {}, "int": {}, "biz": {},
	"info": {}, "name": {}, "pro": {}, "xyz": {}, "online": {}, "site": {}, "shop": {}, "store": {},
	"blog": {}, "news": {}, "media": {}, "art": {}, "film": {}, "game": {}, "games": {}, "tech": {},
	"app": {}, "dev": {}, "ai": {}, "cloud": {}, "io": {}, "co": {}, "me": {}, "tv": {}, "ly": {},
	"to": {}, "fm": {}, "wiki": {}, "help": {}, "us": {}, "uk": {}, "ca": {}, "au": {}, "de": {},
	"fr": {}, "jp": {}, "cn": {}, "ru": {}, "br": {}, "in": {}, "cl": {}, "mx": {}, "es": {},
	"it": {}, "nl": {}, "se": {}, "no": {}, "fi": {}, "dk": {}, "pl": {}, "be": {}, "ch": {},
	"at": {}, "nz": {}, "za": {}, "sg": {}, "hk": {}, "kr": {}, "id": {}, "my": {}, "ph": {},
	"th": {}, "vn": {}, "il": {}, "sa": {}, "ae": {}, "tr": {}, "eg": {}, "ar": {}, "pe": {},
	"ve": {}, "pk": {}, "ng": {}, "ke": {}, "tz": {}, "ro": {},

	// language subdomains (only those not already present above)
	"en": {}, "pt": {}, "zh": {}, "ja": {}, "ko": {}, "sv": {}, "da": {}, "el": {}, "cs": {},
	"hu": {}, "he": {}, "ms": {}, "hi": {}, "bn": {}, "ur": {}, "vi": {},

	// major brands
	"google": {}, "facebook": {}, "instagram": {}, "twitter": {}, "tiktok": {}, "linkedin": {},
	"youtube": {}, "reddit": {}, "wikipedia": {}, "yahoo": {}, "bing": {}, "microsoft": {}, "apple": {},
	"amazon": {}, "ebay": {}, "netflix": {}, "hulu": {}, "spotify": {}, "pinterest": {}, "snapchat": {},
	"discord": {}, "steam": {}, "github": {}, "gitlab": {}, "bitbucket": {}, "twitch": {}, "paypal": {},
	"stripe": {}, "wordpress": {}, "tumblr": {}, "medium": {}, "quora": {}, "stackoverflow": {},
	"dropbox": {}, "icloud": {}, "adobe": {}, "salesforce": {}, "slack": {}, "zoom": {}, "airbnb": {},
	"uber": {}, "lyft": {}, "doordash": {}, "tesla": {}, "openai": {}, "nvidia": {}, "amd": {},
	"intel": {}, "samsung": {}, "huawei": {}, "xiaomi": {}, "sony": {}, "bbc": {}, "cnn": {},
	"nytimes": {}, "forbes": {}, "bloomberg": {}, "wsj": {}, "reuters": {},
}
