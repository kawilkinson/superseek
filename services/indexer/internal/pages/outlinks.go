package pages

type Outlinks struct {
	ID    string
	Links map[string]struct{}
}

func (ol *Outlinks) ToMap() map[string]interface{} {
	outlinks := make([]string, 0, len(ol.Links))
	for link := range ol.Links {
		outlinks = append(outlinks, link)
	}

	return map[string]interface{}{
		"_id":   ol.ID,
		"links": outlinks,
	}
}
