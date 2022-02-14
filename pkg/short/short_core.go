package short

import (
	"shorturl/app/models/short"
	"shorturl/pkg/app"
	"shorturl/pkg/hash"
	"shorturl/pkg/helpers"
	"shorturl/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Generate(longUrl string) (string, bool) {
	// 填充协议头
	longUrl = helpers.FillHttpScheme(longUrl)

	// 哈希拿到短的
	hash := hash.UrlShortHash(longUrl)

	// 查缓存
	shortModel := getCache(hash)
	if shortModel.ID == 0 {
		// 查询db
		shortModel = short.GetByShort(hash)
	}
	if shortModel.ID != 0 {
		if longUrl != shortModel.Long {
			longUrl = longUrl + "#" + helpers.CurrentTimestampStr()
			return Generate(longUrl)
		} else {
			return getShortUrl(shortModel), true
		}
	}

	shortModel = short.CreateShort(hash, longUrl, "0")
	if shortModel.ID == 0 {
		logger.ErrorJSON("short_core", "Generate", gin.H{"hash": hash, "long": longUrl})
		return "", false
	}
	setCache(shortModel)

	return getShortUrl(shortModel), true
}

// 短链接换长链接
func Get(shortStr string) (string, bool) {
	model := short.GetByShort(shortStr)
	if model.ID == 0 {
		return "", false
	}
	setCache(model)
	return model.Long, true
}

func getShortUrl(model short.Short) string {
	return app.URL("/" + model.Short)
}
