package statistic

import (
	"shorturl/app/models"
	"shorturl/pkg/database"
)

type Statistic struct {
	models.BaseModel

	Host  string `gorm:"type:varchar(32);index;not null" json:"host"`
	Short string `gorm:"type:varchar(32);not null;unique" json:"short"`
	Long  string `gorm:"type:text;not null" json:"long"`
	Counter int    `gorm:"type:int unsigned;not null" json:"counter"`

	models.CommonTimestampsField
}

func (model *Statistic) Create() {
	database.DB.Create(&model)
}

func (model *Statistic) Save() int64 {
	result := database.DB.Save(&model)
	return result.RowsAffected
}
