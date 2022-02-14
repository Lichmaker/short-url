package short

import "shorturl/pkg/database"

func GetByShort(short string) Short {
	var model Short
	database.DB.Model(&model).Where("short = ?", short).First(&model)
	return model
}

func CreateShort(short string, long string, expiredAt string) Short {
	model := Short{
		Short:     short,
		Long:      long,
		ExpiredAt: expiredAt,
	}
	model.Create()
	return model
}
