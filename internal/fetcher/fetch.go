package fetcher

import "fmt"

type Fetcher interface {
	Check(url string) bool
	Fetch(url, userAgent string) (string, *SubInfo, error)
}

type SubInfo struct {
	Upload   int64
	Download int64
	Total    int64
	Expire   int64
}

var fetchers = make([]Fetcher, 0)

func init() {
	fetchers = append(fetchers, NewFile())
	fetchers = append(fetchers, NewRemote())
}

func FetchSubscription(url string, userAgent string) (string, *SubInfo, error) {
	for _, fetcher := range fetchers {
		if fetcher.Check(url) {
			return fetcher.Fetch(url, userAgent)
		}
	}
	return "", nil, fmt.Errorf("unsupported url: %s", url)
}
