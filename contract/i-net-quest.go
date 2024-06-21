package contract

import "time"

type INetRequest interface {
	GetBodyJSON() string
	GetHeaders() map[string]string
	GetRoute() string
	GetTimeout() time.Duration
	GetURL() string
}
