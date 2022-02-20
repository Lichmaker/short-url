package statistic

import (
	"shorturl/pkg/database"
	"time"
)

func GetByShort(shortStr string) Statistic {
	var model Statistic
	database.DB.Model(&model).Where("short = ?", shortStr).First(&model)
	return model
}

func IncreaseCounterByShort(shortStr string) {
	database.DB.Exec("UPDATE statistics SET counter = counter + 1, updated_at = ? WHERE short = ?", time.Now(), shortStr)
}
