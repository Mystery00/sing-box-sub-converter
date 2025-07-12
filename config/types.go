package config

type Subscription struct {
	URL       string `json:"url" mapstructure:"url"`
	Tag       string `json:"tag" mapstructure:"tag"`
	Prefix    string `json:"prefix" mapstructure:"prefix"`
	UserAgent string `json:"userAgent" mapstructure:"userAgent"`
}

type ProvidersGlobalConfig struct {
	Subscribes      []Subscription `json:"subscribes" mapstructure:"subscribes"`
	Prefix          bool           `json:"prefix" mapstructure:"prefix"`
	Emoji           bool           `json:"emoji" mapstructure:"emoji"`
	ExcludeProtocol string         `json:"exclude_protocol" mapstructure:"exclude_protocol"`
	ShowSubInNodes  bool           `json:"show_sub_in_nodes" mapstructure:"show_sub_in_nodes"`
}
