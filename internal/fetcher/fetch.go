package fetcher

import "fmt"

type Fetcher interface {
	Check(url string) bool
	Fetch(url, userAgent string) (string, error)
}

var fetchers = make([]Fetcher, 0)

func init() {
	fetchers = append(fetchers, NewFile())
	fetchers = append(fetchers, NewRemote())
}

func FetchSubscription(url string, userAgent string) (string, error) {
	for _, fetcher := range fetchers {
		if fetcher.Check(url) {
			return fetcher.Fetch(url, userAgent)
		}
	}
	return "", fmt.Errorf("unsupported url: %s", url)
}
