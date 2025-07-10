package indexerutil

import "time"

const (
	Timeout       = 4 * time.Second
	MaxIndexWords = 1000
	WordImagesOpThreshold = 500
	ImageOpThreshold = 100
	ImgMinWidth   = 100
	ImgMinHeight  = 500

	WordCollection       = "words"
	WordImagesCollection = "word_images"
	MetadataCollection   = "metadata"
	ImageCollection      = "images"

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

	BacklinksPrefix = "backlinks"
	OutlinksPrefix  = "outlinks"
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

// stop words grabbed from here https://gist.github.com/sebleier/554280
var StopWordsSet = map[string]struct{}{
	"i": {}, "me": {}, "my": {}, "myself": {}, "we": {}, "our": {},
	"ours": {}, "ourselves": {}, "you": {}, "your": {}, "yours": {},
	"yourself": {}, "yourselves": {}, "he": {}, "him": {}, "his": {},
	"himself": {}, "she": {}, "her": {}, "hers": {}, "herself": {},
	"it": {}, "its": {}, "itself": {}, "they": {}, "them": {}, "their": {},
	"theirs": {}, "themselves": {}, "what": {}, "which": {}, "who": {},
	"whom": {}, "this": {}, "that": {}, "these": {}, "those": {},
	"am": {}, "is": {}, "are": {}, "was": {}, "were": {}, "be": {},
	"been": {}, "being": {}, "have": {}, "has": {}, "had": {},
	"having": {}, "do": {}, "does": {}, "did": {}, "doing": {}, "a": {},
	"an": {}, "the": {}, "and": {}, "but": {}, "if": {}, "or": {},
	"because": {}, "as": {}, "until": {}, "while": {}, "of": {}, "at": {},
	"by": {}, "for": {}, "with": {}, "about": {}, "against": {},
	"between": {}, "into": {}, "through": {}, "during": {}, "before": {},
	"after": {}, "above": {}, "below": {}, "to": {}, "from": {}, "up": {},
	"down": {}, "in": {}, "out": {}, "on": {}, "off": {}, "over": {},
	"under": {}, "again": {}, "further": {}, "then": {}, "once": {},
	"here": {}, "there": {}, "when": {}, "where": {}, "why": {}, "how": {},
	"all": {}, "any": {}, "both": {}, "each": {}, "few": {}, "more": {},
	"most": {}, "other": {}, "some": {}, "such": {}, "no": {}, "nor": {},
	"not": {}, "only": {}, "own": {}, "same": {}, "so": {}, "than": {},
	"too": {}, "very": {}, "s": {}, "t": {}, "can": {}, "will": {},
	"just": {}, "don": {}, "should": {}, "now": {},
}
