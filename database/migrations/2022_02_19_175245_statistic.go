package migrations

import (
	"database/sql"
	"shorturl/app/models/statistic"
	"shorturl/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.AutoMigrate(&statistic.Statistic{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.DropTable(&statistic.Statistic{})
	}

	migrate.Add("2022_02_19_175245_statistic", up, down)
}
