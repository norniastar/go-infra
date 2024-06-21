package rediskey

type Value string

func (m Value) String() string {
	return string(m)
}

const (
	ConfigCache Value = "cache-config"
)
