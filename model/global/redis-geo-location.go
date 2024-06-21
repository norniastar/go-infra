package global

type RedisGeoLocation struct {
	RedisGeoPosition

	Distance float64
	Hash     int64
	Member   string
}
