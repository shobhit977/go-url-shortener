package models

type UrlInfo struct {
	Url      string `json:"url"`
	ShortUrl string `json:"shorturl"`
	Domain   string `json:"domain"`
}
type MetricsResponse struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}

type UrlShortenerRequest struct {
	Url string `json:"url"`
}
