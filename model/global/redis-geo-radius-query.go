package global

type RedisGeoRadiusQuery struct {
	Count     int
	Radius    float64
	Sort      string
	Unit      string
	WithCoord bool
	WithDist  bool
	WithHash  bool
}
