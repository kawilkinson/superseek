package pages

type PageNode struct {
	NormalizedURL      string
	NormalizedLinkURLs map[string]struct{}
}

func CreatePageNode(normalizedURL string) *PageNode {
	return &PageNode{
		NormalizedURL:      normalizedURL,
		NormalizedLinkURLs: make(map[string]struct{}),
	}
}

func (pn *PageNode) AppendLink(newNormalizedLink string) {
	if pn.NormalizedLinkURLs == nil {
		pn.NormalizedLinkURLs = make(map[string]struct{})
	}

	pn.NormalizedLinkURLs[newNormalizedLink] = struct{}{}
}

func (pn *PageNode) GetLinks() []string {
	var links []string
	for link := range pn.NormalizedLinkURLs {
		links = append(links, link)
	}

	return links
}
