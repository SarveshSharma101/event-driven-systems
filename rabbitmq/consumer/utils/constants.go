package utils

type ExchangeType string

const (
	DIRECT  ExchangeType = "direct"
	TOPIC   ExchangeType = "topic"
	FANOUT  ExchangeType = "fanout"
	HEADERS ExchangeType = "headers"

	URL string = "url"
)
