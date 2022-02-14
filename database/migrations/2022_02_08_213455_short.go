package migrations

import (
	"database/sql"
	"shorturl/app/models/short"
	"shorturl/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.AutoMigrate(&short.Short{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.DropTable(&short.Short{})
	}

	migrate.Add("2022_02_08_213455_short", up, down)
}
