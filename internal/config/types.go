package config

type Subscription struct {
	URL       string `json:"url"`
	Tag       string `json:"tag"`
	Prefix    string `json:"prefix"`
	UserAgent string `json:"userAgent"`
}

type ProvidersGlobalConfig struct {
	Subscribes      []Subscription `json:"subscribes"`
	Prefix          bool           `json:"prefix"`
	Emoji           bool           `json:"emoji"`
	ExcludeProtocol string         `json:"exclude_protocol"`
}
