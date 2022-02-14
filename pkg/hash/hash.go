package hash

import (
	"shorturl/pkg/helpers"
	"shorturl/pkg/logger"

	"github.com/spaolacci/murmur3"
	"golang.org/x/crypto/bcrypt"
)

func BcryptHash(psw string) string {
	hashByte, err := bcrypt.GenerateFromPassword([]byte(psw), bcrypt.DefaultCost)
	logger.LogIf(err)

	return string(hashByte)
}

func BcryptCheck(psw, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(psw))
	return err == nil
}

func BcryptIsHashed(s string) bool {
	return len(s) == 60
}

func UrlShortHash(url string) string {
	// 使用32位的murmur3算法，计算出32位数字
	murmur3Hash := murmur3.Sum32([]byte(url))

	// 使用uint的base62，得出字符串
	return helpers.ConvNumToBase62(uint(murmur3Hash))
}
