package goredissvc

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"github.com/norniastar/infra-core/contract"
	"github.com/norniastar/infra-core/model/global"
	"github.com/samber/lo"
	"reflect"
	"strconv"
	"sync"
	"time"
)

var redisAdapterMutex sync.Mutex

type redisAdapter struct {
	client         redis.Cmdable
	clusterOptions *redis.ClusterOptions
	ctx            context.Context
	options        *redis.Options
}

func (r *redisAdapter) BitCount(key string, start, end int64) (int64, error) {
	return r.getClient().BitCount(key, &redis.BitCount{
		End:   end,
		Start: start,
	}).Result()
}

func (r *redisAdapter) BitOp(op, destKey string, keys ...string) (bool, error) {
	var res int64
	var err error
	switch op {
	case "and":
		res, err = r.getClient().BitOpAnd(destKey, keys...).Result()
	case "not":
		res, err = r.getClient().BitOpNot(destKey, keys[0]).Result()
	case "or":
		res, err = r.getClient().BitOpOr(destKey, keys...).Result()
	default:
		res, err = r.getClient().BitOpXor(destKey, keys...).Result()
	}

	if err != nil {
		return false, err
	}

	temp := false
	if res == 1 {
		temp = true
	}
	return temp, err
}

func (r *redisAdapter) BitPos(key string, bit bool, start, end int64) (int64, error) {
	v := int64(0)
	if bit {
		v = 1
	}

	return r.getClient().BitPos(key, v, start, end).Result()
}

func (r *redisAdapter) BLPop(timeout time.Duration, keys ...string) ([]string, error) {
	return r.getClient().BLPop(timeout, keys...).Result()
}

func (r *redisAdapter) BRPop(timeout time.Duration, keys ...string) ([]string, error) {
	return r.getClient().BRPop(timeout, keys...).Result()
}

func (r *redisAdapter) Close() error {
	if c, ok := r.getClient().(*redis.Client); ok {
		return c.Close()
	}

	return r.getClient().(*redis.ClusterClient).Close()
}

func (r *redisAdapter) Decr(key string) (int64, error) {
	return r.getClient().Decr(key).Result()
}

func (r *redisAdapter) DecrBy(key string, decrement int64) (int64, error) {
	return r.getClient().DecrBy(key, decrement).Result()
}

func (r *redisAdapter) Del(keys ...string) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	return r.getClient().Del(keys...).Result()
}

func (r *redisAdapter) Exists(str string) (bool, error) {
	res := r.getClient().Exists(str)
	return res.Val() == 1, res.Err()
}

func (r *redisAdapter) ExpireAt(key string, time time.Time) (bool, error) {
	return r.getClient().ExpireAt(key, time).Result()
}

func (r *redisAdapter) Expires(key string, seconds time.Duration) (bool, error) {
	return r.getClient().Expire(key, seconds).Result()
}

func (r *redisAdapter) GeoAdd(key string, locations ...global.RedisGeoLocation) (int64, error) {
	temp := lo.Map(locations, func(r global.RedisGeoLocation, _ int) *redis.GeoLocation {
		return &redis.GeoLocation{
			Latitude:  r.Latitude,
			Longitude: r.Longitude,
			Name:      r.Member,
		}
	})

	return r.getClient().GeoAdd(key, temp...).Result()
}

func (r *redisAdapter) GeoDist(key string, member1, member2, unit string) (float64, error) {
	return r.getClient().GeoDist(key, member1, member2, unit).Result()
}

func (r *redisAdapter) GeoPos(key string, members ...string) ([]*global.RedisGeoPosition, error) {
	res, err := r.getClient().GeoPos(key, members...).Result()
	if err != nil {
		return nil, err
	}

	return lo.FilterMap(res, func(r *redis.GeoPos, _ int) (*global.RedisGeoPosition, bool) {
		if r == nil {
			return &global.RedisGeoPosition{}, false
		}

		return &global.RedisGeoPosition{
			Latitude:  r.Latitude,
			Longitude: r.Longitude,
		}, true
	}), nil
}

func (r *redisAdapter) GeoRadius(key string, longitude, latitude float64, query global.RedisGeoRadiusQuery) ([]global.RedisGeoLocation, error) {
	res, err := r.getClient().GeoRadius(key, longitude, latitude, &redis.GeoRadiusQuery{
		Count:       query.Count,
		Radius:      query.Radius,
		Sort:        query.Sort,
		Unit:        query.Unit,
		WithCoord:   query.WithCoord,
		WithDist:    query.WithDist,
		WithGeoHash: query.WithHash,
	}).Result()
	if err != nil {
		return nil, err
	}

	return lo.Map(res, func(r redis.GeoLocation, _ int) global.RedisGeoLocation {
		return global.RedisGeoLocation{
			RedisGeoPosition: global.RedisGeoPosition{
				Latitude:  r.Latitude,
				Longitude: r.Longitude,
			},
			Distance: r.Dist,
			Hash:     r.GeoHash,
			Member:   r.Name,
		}
	}), nil
}

func (r *redisAdapter) GeoRadiusByMember(key string, member string, query global.RedisGeoRadiusQuery) ([]global.RedisGeoLocation, error) {
	res, err := r.getClient().GeoRadiusByMember(key, member, &redis.GeoRadiusQuery{
		Count:       query.Count,
		Radius:      query.Radius,
		Sort:        query.Sort,
		Unit:        query.Unit,
		WithCoord:   query.WithCoord,
		WithDist:    query.WithDist,
		WithGeoHash: query.WithHash,
	}).Result()
	if err != nil {
		return nil, err
	}

	return lo.Map(res, func(item redis.GeoLocation, _ int) global.RedisGeoLocation {
		return global.RedisGeoLocation{
			RedisGeoPosition: global.RedisGeoPosition{
				Latitude:  item.Latitude,
				Longitude: item.Longitude,
			},
			Distance: item.Dist,
			Hash:     item.GeoHash,
			Member:   item.Name,
		}
	}), nil
}

func (r *redisAdapter) Get(str string) (string, error) {
	res, err := r.getClient().Get(str).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return "", nil
	}

	return res, err
}

func (r *redisAdapter) GetBit(key string, offset int64) (bool, error) {
	res, err := r.getClient().GetBit(key, offset).Result()
	if err != nil {
		return false, err
	}

	return res == 1, nil
}

func (r *redisAdapter) HDel(key string, fields ...string) (int64, error) {
	return r.getClient().HDel(key, fields...).Result()
}

func (r *redisAdapter) HExists(key, field string) (bool, error) {
	return r.getClient().HExists(key, field).Result()
}

func (r *redisAdapter) HGet(key, field string) (string, error) {
	res, err := r.getClient().HGet(key, field).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return "", nil
	}

	return res, err
}

func (r *redisAdapter) HGetAll(key string) (map[string]string, error) {
	return r.getClient().HGetAll(key).Result()
}

func (r *redisAdapter) HIncrBy(key, field string, increment int64) (int64, error) {
	return r.getClient().HIncrBy(key, field, increment).Result()
}

func (r *redisAdapter) HIncrByFloat(key, field string, increment float64) (float64, error) {
	return r.getClient().HIncrByFloat(key, field, increment).Result()
}

func (r *redisAdapter) HKeys(key string) ([]string, error) {
	return r.getClient().HKeys(key).Result()
}

func (r *redisAdapter) HLen(key string) (int64, error) {
	return r.getClient().HLen(key).Result()
}

func (r *redisAdapter) HMGet(key string, fields ...string) ([]string, error) {
	res, err := r.getClient().HMGet(key, fields...).Result()
	if err != nil {
		return nil, err
	}

	return lo.Map(res, func(r any, _ int) string {
		return r.(string)
	}), nil
}

func (r *redisAdapter) HMSet(key string, values map[string]interface{}) error {
	_, err := r.getClient().HMSet(key, values).Result()
	return err
}

func (r *redisAdapter) HScan(key string, cursor uint64, match string, count int64) (map[string]string, uint64, error) {
	res, cursor, err := r.getClient().HScan(key, cursor, match, count).Result()
	if err != nil {
		return nil, cursor, err
	}

	chunk := lo.Chunk(res, 2)
	temp := lo.SliceToMap(chunk, func(r []string) (string, string) {
		return r[0], r[1]
	})
	return temp, cursor, nil
}

func (r *redisAdapter) HSet(key, field, value string) (bool, error) {
	return r.getClient().HSet(key, field, value).Result()
}

func (r *redisAdapter) HSetNX(key, field, value string) (bool, error) {
	return r.getClient().HSetNX(key, field, value).Result()
}

func (r *redisAdapter) HStrLen(key, field string) (int64, error) {
	res, err := r.HGet(key, field)
	if err != nil {
		return 0, err
	}

	return int64(len(res)), nil
}

func (r *redisAdapter) HVals(key string) ([]string, error) {
	return r.getClient().HVals(key).Result()
}

func (r *redisAdapter) Incr(key string) (int64, error) {
	return r.getClient().Incr(key).Result()
}

func (r *redisAdapter) IncrBy(key string, increment int64) (int64, error) {
	return r.getClient().IncrBy(key, increment).Result()
}

func (r *redisAdapter) LIndex(key string, index int64) (string, error) {
	return r.getClient().LIndex(key, index).Result()
}

func (r *redisAdapter) LLen(key string) (int64, error) {
	return r.getClient().LLen(key).Result()
}

func (r *redisAdapter) LPop(key string) (string, error) {
	return r.getClient().LPop(key).Result()
}

func (r *redisAdapter) LPush(key string, values ...string) (int64, error) {
	if len(values) == 0 {
		return 0, nil
	}

	temp := lo.Map(values, func(r string, _ int) any {
		return r
	})
	return r.getClient().LPush(key, temp...).Result()
}

func (r *redisAdapter) LPushX(key string, value string) (int64, error) {
	return r.getClient().LPushX(key, value).Result()
}

func (r *redisAdapter) LRange(key string, start, stop int64) ([]string, error) {
	return r.getClient().LRange(key, start, stop).Result()
}

func (r *redisAdapter) LRem(key string, count int64, value string) (int64, error) {
	return r.getClient().LRem(key, count, value).Result()
}

func (r *redisAdapter) LSet(key string, index int64, value string) (bool, error) {
	res, err := r.getClient().LSet(key, index, value).Result()
	if err != nil {
		return false, err
	}

	return res == "OK", nil
}

func (r *redisAdapter) LTrim(key string, start, stop int64) (bool, error) {
	res, err := r.getClient().LTrim(key, start, stop).Result()
	if err != nil {
		return false, err
	}

	return res == "OK", nil
}

func (r *redisAdapter) RPop(key string) (string, error) {
	return r.getClient().RPop(key).Result()
}

func (r *redisAdapter) RPush(key string, values ...string) (int64, error) {
	if len(values) == 0 {
		return 0, nil
	}

	temp := lo.Map(values, func(r string, _ int) any {
		return r
	})
	return r.getClient().RPush(key, temp...).Result()
}

func (r *redisAdapter) RPushX(key string, value string) (int64, error) {
	return r.getClient().RPushX(key, value).Result()
}

func (r *redisAdapter) SAdd(key string, members ...string) (int64, error) {
	temp := lo.Map(members, func(r string, _ int) any {
		return r
	})
	return r.getClient().SAdd(key, temp...).Result()
}

func (r *redisAdapter) SCard(key string) (int64, error) {
	return r.getClient().SCard(key).Result()
}

func (r *redisAdapter) Set(key, value string, extraArgs ...interface{}) (ok bool, err error) {
	var res string
	if len(extraArgs) == 0 {
		res, err = r.getClient().Set(key, value, 0).Result()
		ok = res == "OK"
	} else if len(extraArgs) == 1 {
		if extraArgs[0] == "nx" {
			ok, err = r.getClient().SetNX(key, value, 0).Result()
		} else if extraArgs[0] == "xx" {
			ok, err = r.getClient().SetXX(key, value, 0).Result()
		} else {
			panic("redis set 参数有误")
		}
	} else if len(extraArgs) == 2 {
		t := reflect.ValueOf(extraArgs[1]).Int()
		var expires time.Duration
		if extraArgs[0] == "ex" {
			expires = time.Duration(t) * time.Second
		} else if extraArgs[0] == "px" {
			expires = time.Duration(t) * time.Millisecond
		} else {
			panic("redis set 参数有误")
		}
		res, err = r.getClient().Set(key, value, expires).Result()
		ok = res == "OK"
	} else if len(extraArgs) == 3 {
		t := reflect.ValueOf(extraArgs[1]).Int()
		var expires time.Duration
		if extraArgs[0] == "ex" {
			expires = time.Duration(t) * time.Second
		} else if extraArgs[0] == "px" {
			expires = time.Duration(t) * time.Millisecond
		} else {
			panic("redis set 参数有误")
		}

		if extraArgs[2] == "nx" {
			ok, err = r.getClient().SetNX(key, value, expires).Result()
		} else if extraArgs[2] == "xx" {
			ok, err = r.getClient().SetXX(key, value, expires).Result()
		} else {
			panic("redis set 参数有误")
		}
	} else {
		panic("redis set 参数过多")
	}
	return
}

func (r *redisAdapter) SetBit(key string, offset int64, value bool) (bool, error) {
	temp := 0
	if value {
		temp = 1
	}
	res, err := r.getClient().SetBit(key, offset, temp).Result()
	if err != nil {
		return false, err
	}

	return res == 1, nil
}

func (r *redisAdapter) SIsMember(key, member string) (bool, error) {
	return r.getClient().SIsMember(key, member).Result()
}

func (r *redisAdapter) SMembers(key string) ([]string, error) {
	return r.getClient().SMembers(key).Result()
}

func (r *redisAdapter) SPop(key string) (string, error) {
	res, err := r.getClient().SPop(key).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return "", nil
	}
	return res, nil
}

func (r *redisAdapter) Time() (time.Time, error) {
	return r.getClient().Time().Result()
}

func (r *redisAdapter) TTL(key string) (time.Duration, error) {
	return r.getClient().TTL(key).Result()
}

func (r *redisAdapter) WithContext(ctx context.Context) reflect.Value {
	return reflect.ValueOf(&redisAdapter{
		client: r.getClient(),
		ctx:    ctx,
	})
}

func (r *redisAdapter) ZAdd(key string, members ...global.RedisZMember) (int64, error) {
	if len(members) == 0 {
		return 0, nil
	}

	temp := lo.Map(members, func(r global.RedisZMember, _ int) redis.Z {
		return redis.Z{
			Member: r.Member,
			Score:  r.Score,
		}
	})
	return r.getClient().ZAdd(key, temp...).Result()
}

func (r *redisAdapter) ZCard(key string) (int64, error) {
	return r.getClient().ZCard(key).Result()
}

func (r *redisAdapter) ZCount(key string, min, max float64) (int64, error) {
	return r.getClient().ZCount(

		key,
		strconv.FormatFloat(min, 'E', -1, 64),
		strconv.FormatFloat(max, 'E', -1, 64),
	).Result()
}

func (r *redisAdapter) ZIncrBy(key string, increment float64, member string) (float64, error) {
	return r.getClient().ZIncrBy(key, increment, member).Result()
}

func (r *redisAdapter) ZRange(key string, start, stop int64, withScores bool) ([]global.RedisZMember, error) {
	var members []global.RedisZMember
	if withScores {
		res, err := r.getClient().ZRangeWithScores(key, start, stop).Result()
		if err != nil {
			return nil, err
		} else if len(res) > 0 {

			members = lo.Map(res, func(r redis.Z, _ int) global.RedisZMember {
				return global.RedisZMember{
					Member: r.Member.(string),
					Score:  r.Score,
				}
			})
		}
	} else {
		res, err := r.getClient().ZRange(key, start, stop).Result()
		if err != nil {
			return nil, err
		} else if len(res) > 0 {
			members = lo.Map(res, func(r string, _ int) global.RedisZMember {
				return global.RedisZMember{
					Member: r,
				}
			})
		}
	}

	return members, nil
}

func (r *redisAdapter) ZRangeByScore(key string, min, max string, opt global.RedisZRangeByScore) ([]global.RedisZMember, error) {
	var members []global.RedisZMember
	if opt.WithScores {
		res, err := r.getClient().ZRangeByScoreWithScores(key, redis.ZRangeBy{
			Count:  opt.Count,
			Max:    max,
			Min:    min,
			Offset: opt.Offset,
		}).Result()
		if err != nil {
			return nil, err
		} else if len(res) > 0 {
			members = lo.Map(res, func(r redis.Z, _ int) global.RedisZMember {
				return global.RedisZMember{
					Member: r.Member.(string),
					Score:  r.Score,
				}
			})
		}
	} else {
		res, err := r.getClient().ZRangeByScore(key, redis.ZRangeBy{
			Count:  opt.Count,
			Max:    max,
			Min:    min,
			Offset: opt.Offset,
		}).Result()
		if err != nil {
			return nil, err
		} else if len(res) > 0 {
			members = lo.Map(res, func(r string, _ int) global.RedisZMember {
				return global.RedisZMember{
					Member: r,
				}
			})
		}
	}

	return members, nil
}

func (r *redisAdapter) ZRank(key, member string) (int64, error) {
	res, err := r.getClient().ZRank(key, member).Result()
	if errors.Is(err, redis.Nil) {
		res = -1
		err = nil
	}
	return res, err
}

func (r *redisAdapter) ZRem(key string, members ...string) (int64, error) {
	temp := lo.Map(members, func(r string, _ int) interface{} {
		return r
	})
	return r.getClient().ZRem(key, temp...).Result()
}

func (r *redisAdapter) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	return r.getClient().ZRemRangeByRank(key, start, stop).Result()
}

func (r *redisAdapter) ZRemRangeByScore(key string, min, max float64) (int64, error) {
	return r.getClient().ZRemRangeByScore(

		key,
		strconv.FormatFloat(min, 'E', -1, 64),
		strconv.FormatFloat(max, 'E', -1, 64),
	).Result()
}

func (r *redisAdapter) ZRevRange(key string, start, stop int64, withScores bool) ([]global.RedisZMember, error) {
	var members []global.RedisZMember
	if withScores {
		res, err := r.getClient().ZRevRangeWithScores(key, start, stop).Result()
		if err != nil {
			return nil, err
		} else if len(res) > 0 {
			members = lo.Map(res, func(r redis.Z, _ int) global.RedisZMember {
				return global.RedisZMember{
					Member: r.Member.(string),
					Score:  r.Score,
				}
			})
		}
	} else {
		res, err := r.getClient().ZRevRange(key, start, stop).Result()
		if err != nil {
			return nil, err
		} else if len(res) > 0 {
			members = lo.Map(res, func(r string, _ int) global.RedisZMember {
				return global.RedisZMember{
					Member: r,
				}
			})
		}
	}

	return members, nil
}

func (r *redisAdapter) ZRevRangeByScore(key string, min, max string, opt global.RedisZRangeByScore) ([]global.RedisZMember, error) {
	var members []global.RedisZMember
	if opt.WithScores {
		res, err := r.getClient().ZRevRangeByScoreWithScores(key, redis.ZRangeBy{
			Count:  opt.Count,
			Max:    max,
			Min:    min,
			Offset: opt.Offset,
		}).Result()
		if err != nil {
			return nil, err
		} else if len(res) > 0 {
			members = lo.Map(res, func(r redis.Z, _ int) global.RedisZMember {
				return global.RedisZMember{
					Member: r.Member.(string),
					Score:  r.Score,
				}
			})
		}
	} else {
		res, err := r.getClient().ZRevRangeByScore(key, redis.ZRangeBy{
			Count:  opt.Count,
			Max:    max,
			Min:    min,
			Offset: opt.Offset,
		}).Result()
		if err != nil {
			return nil, err
		} else if len(res) > 0 {
			members = lo.Map(res, func(r string, _ int) global.RedisZMember {
				return global.RedisZMember{
					Member: r,
				}
			})
		}
	}

	return members, nil
}

func (r *redisAdapter) ZRevRank(key, member string) (int64, error) {
	return r.getClient().ZRevRank(key, member).Result()
}

func (r *redisAdapter) ZScan(key string, cursor uint64, match string, count int64) ([]global.RedisZMember, uint64, error) {
	res, cursor, err := r.getClient().ZScan(key, cursor, match, count).Result()
	if err != nil {
		return nil, cursor, err
	}

	chunk := lo.Chunk(res, 2)
	members := lo.Map(chunk, func(r []string, _ int) global.RedisZMember {
		score, _ := strconv.ParseFloat(r[1], 64)
		return global.RedisZMember{
			Member: r[0],
			Score:  score,
		}
	})
	return members, cursor, nil
}

func (r *redisAdapter) ZScore(key, member string) (float64, error) {
	res, err := r.getClient().ZScore(key, member).Result()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return res, err
}

func (r *redisAdapter) getClient() redis.Cmdable {
	if r.client == nil {
		redisAdapterMutex.Lock()
		defer redisAdapterMutex.Unlock()

		if r.client == nil {
			if r.options != nil {
				r.client = redis.NewClient(r.options).WithContext(r.ctx)
			} else {
				r.client = redis.NewClusterClient(r.clusterOptions).WithContext(r.ctx)
			}
		}
	}
	return r.client
}

func (r *redisAdapter) New(options ...contract.RedisOption) contract.IRedis {
	adapter := &redisAdapter{
		ctx: context.Background(),
	}
	for _, cr := range options {
		cr(adapter)
	}
	return adapter
}
