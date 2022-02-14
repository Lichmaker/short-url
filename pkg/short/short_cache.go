package short

import (
	"shorturl/app/models/short"
	"shorturl/pkg/cache"
	"shorturl/pkg/config"
	"shorturl/pkg/helpers"
	"shorturl/pkg/logger"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
)

// 每生成一个短链接都记录缓存
func setCache(model short.Short) {
	cache.Set(shortStringCacheKey(model.Short), model, 0)

	cache.AddSortedSet(shortSortedSetCacheKey(), model.Short, cast.ToFloat64(helpers.CurrentTimestampStr()))
}

func getCache(shortStr string) short.Short {
	var model short.Short
	cache.GetObject(shortStringCacheKey(shortStr), &model)
	return model
}

func shortStringCacheKey(short string) string {
	return "short_url:" + short
}

func shortSortedSetCacheKey() string {
	return "short_url:my_all_key"
}

func TidyCache() int {
	var count int
	trashData := getTrashCache()
	for _, d := range trashData {
		cache.RemSortedSetMember(shortSortedSetCacheKey(), d.Member.(string))
		cache.Forget(shortStringCacheKey(d.Member.(string)))
		count++
	}
	return count
}

func getTrashCache() []redis.Z {
	// 只保留前1000个key的缓存，其他都丢掉
	trashData, err := cache.RevRangeSortedSet(shortSortedSetCacheKey(), config.GetInt64("short.cache_max"), -1)
	logger.LogIf(err)
	return trashData
}
