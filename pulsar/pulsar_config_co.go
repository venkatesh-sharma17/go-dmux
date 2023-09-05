package pulsar

type PulsarConf struct {
	SubscriptionName string `json:"name"`
	Url              string `json:"url"`
	Topic            string `json:"topic"`
	ForceRestart     bool   `json:"force_restart"`
	ReadNewest       bool   `json:"read_newest"`
	SeekByTime       int64  `json:"seek_by_time"`
	AuthClientId     string `json:"client_id"`
	AuthClientSecret string `json:"auth_client_secret"`
	AuthIssuerURL    string `json:"auth_issuer_url"`
	AuthAudience     string `json:"auth_audience"`
	SubscriptionType string `json:"subscription_type"`
}
