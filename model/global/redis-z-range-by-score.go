package global

type RedisZRangeByScore struct {
	Count, Offset int64
	WithScores    bool
}
