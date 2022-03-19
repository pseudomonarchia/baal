package migrationfile

import (
	"baal/database/migrations"
	"baal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	migrations.SetMigration(&gormigrate.Migration{
		ID: "20220320025032_user_hash",
		Migrate: func(db *gorm.DB) error {
			model := &models.UserHashSchema{}
			return db.Migrator().AutoMigrate(model)
		},
		Rollback: func(db *gorm.DB) error {
			model := &models.UserHashSchema{}
			return db.Migrator().DropTable(model.TableName())
		},
	})
}
