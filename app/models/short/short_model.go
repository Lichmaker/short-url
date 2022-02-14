package short

import (
	"shorturl/app/models"
	"shorturl/pkg/database"
)

type Short struct {
	models.BaseModel

	Short     string `gorm:"type:varchar(32);not null;unique" json:"short"`
	Long      string `gorm:"type:text;not null" json:"long"`
	ExpiredAt string `gorm:"type:int(11) unsigned;index" json:"expried_at"`

	models.CommonTimestampsField
}

func (model *Short) Create() {
	database.DB.Create(&model)
}

func (model *Short) Save() int64 {
	result := database.DB.Save(&model)
	return result.RowsAffected
}
